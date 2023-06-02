mod cli;
mod monitor;
use anyhow::Result;
use futures::future;
use std::{env, process};
use tokio::{signal, spawn};

use crate::cli::cli;
use crate::monitor::configure_monitors;
use roselite_config::Configuration;
use roselite_server::config::ServerConfig;

#[tokio::main]
async fn main() -> Result<()> {
    let matches = cli().get_matches();

    // There are 3 ways to set a configuration file path:
    //   - `CONFIGURATION_FILE_PATH` environment variable
    //   - `-c` or `--config` CLI arguments
    //   - by a default value of `conf.toml`.
    let mut configuration_file_path: String =
        env::var("CONFIGURATION_FILE_PATH").unwrap_or(String::from("conf.toml"));
    if let Some(conf) = matches.get_one::<String>("config") {
        configuration_file_path = conf.to_string();
    }

    let configuration: Configuration = Configuration::from_file(&configuration_file_path)?;

    // Configure Sentry as default error monitoring
    let sentry_dsn: String =
        env::var("SENTRY_DSN").unwrap_or(match configuration.error_reporting {
            Some(error_reporting) => error_reporting.sentry_dsn,
            None => String::new(),
        });
    let _guard = sentry::init(sentry_dsn);

    let mut handles = vec![];

    match matches.subcommand() {
        Some(("server", _)) => {
            let monitor_handles = configure_monitors(configuration.monitors);
            handles.extend(monitor_handles);

            if let Some(server) = configuration.server {
                handles.push(spawn(async {
                    // Start server
                    println!("HTTP server is starting at {}", server.listen_address);
                    roselite_server::run(ServerConfig {
                        address: server.listen_address,
                        // TODO: upstream support soon
                        upstream_kuma: server.upstream_kuma,
                    })
                    .await
                    .unwrap()
                }));
            } else {
                // If configuration.server is None, we can't continue do anything
                eprintln!("configuration.server must be filled, either way, don't run it on 'server' mode");
                process::exit(10);
            }
        }
        _ => {
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
