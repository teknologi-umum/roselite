mod config;
mod task_runner;

use std::{env, process};
use tokio::{signal, spawn};
use tokio::time::{Instant, Duration, sleep};
use anyhow::Result;
use futures::future;
use crate::config::Configuration;
use crate::task_runner::perform_task;

#[tokio::main]
async fn main() -> Result<()> {
    let configuration_file_path = env::var("CONFIGURATION_FILE_PATH").unwrap_or(String::from("conf.toml"));
    let configuration = Configuration::from_file(&configuration_file_path)?;

    let mut handles = vec![];
    for monitor in configuration.monitors {
        println!("Starting monitor for {}", monitor.monitor_url);
        handles.push(spawn(async move {
            loop {
                let current_time = Instant::now();
                if let Err(err) = perform_task(monitor.clone()).await {
                    // Do nothing of this error
                    println!("{:?}", err);
                }

                let elapsed = current_time.elapsed();
                if 60 - elapsed.as_secs() > 0 {
                    sleep(Duration::from_secs(60 - elapsed.as_secs())).await;
                    continue;
                }

                // Immediately continue
                continue;
            }
        }));
    }

    match signal::ctrl_c().await {
        Ok(()) => {
            println!("Received a shutdown signal, exiting...");
            process::exit(0);
        },
        Err(err) => {
            eprintln!("Unable to listen for shutdown signal: {}", err);
        }
    }

    future::join_all(handles).await;

    Ok(())
}
