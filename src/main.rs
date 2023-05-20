use std::{env, process};

use anyhow::Result;
use futures::future;
use tokio::time::{sleep, Duration, Instant};
use tokio::{signal, spawn};

use crate::config::{Configuration, Monitor};
use crate::task_runner::perform_task;

mod config;
mod task_runner;

#[tokio::main]
async fn main() -> Result<()> {
    let configuration_file_path: String =
        env::var("CONFIGURATION_FILE_PATH").unwrap_or(String::from("conf.toml"));
    let configuration: Configuration = Configuration::from_file(&configuration_file_path)?;

    let mut handles = vec![];
    for monitor in configuration.monitors {
        println!("Starting monitor for {}", monitor.monitor_url);
        handles.push(spawn(async move {
            let cloned_monitor: Monitor = monitor.clone();
            loop {
                let current_time = Instant::now();
                if let Err(err) = perform_task(monitor.clone()).await {
                    // Do nothing of this error
                    eprintln!("Unexpected error during performing task: {}", err);
                }

                let elapsed = current_time.elapsed();
                if 60 - elapsed.as_secs() > 0 {
                    let sleeping_duration = Duration::from_secs(60 - elapsed.as_secs());
                    println!(
                        "Monitor for {0} will be sleeping for {1} seconds",
                        cloned_monitor.monitor_url,
                        sleeping_duration.as_secs()
                    );
                    sleep(sleeping_duration).await;
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
        }
        Err(err) => {
            eprintln!("Unable to listen for shutdown signal: {}", err);
        }
    }

    // This is here just because we don't want the application to immediately exits.
    future::join_all(handles).await;

    Ok(())
}
