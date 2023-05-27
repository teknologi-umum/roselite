mod cli;
mod monitor;
use anyhow::Result;
use futures::future;
use std::{env, process};
use tokio::{signal, spawn};

use crate::cli::cli;
use crate::monitor::configure_monitors;
use roselite_config::Configuration;

#[tokio::main]
async fn main() -> Result<()> {
    let sentry_dsn: String = env::var("SENTRY_DSN").unwrap_or(String::from(""));
    let _guard = sentry::init(sentry_dsn);

    let configuration_file_path: String =
        env::var("CONFIGURATION_FILE_PATH").unwrap_or(String::from("conf.toml"));

    let mut handles = vec![];
    let matches = cli().get_matches();
    match matches.subcommand() {
        Some(("server", _)) => {
            let configuration: Configuration = Configuration::from_file(&configuration_file_path)?;
            let monitor_handles = configure_monitors(configuration.monitors);
            handles.extend(monitor_handles);

            // Start server
            // TODO: configure port and host
            handles.push(spawn(async {
                roselite_server::run("127.0.0.1:8321".into())
                    .await
                    .expect("TODO: panic message");
            }));
        }
        _ => {
            let configuration: Configuration = Configuration::from_file(&configuration_file_path)?;
            let monitor_handles = configure_monitors(configuration.monitors);
            handles.extend(monitor_handles);
        }
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
