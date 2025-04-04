use std::time::Duration;

use sentry::integrations::anyhow::capture_anyhow;
use tokio::spawn;
use tokio::task::JoinHandle;
use tokio::time::{Instant, sleep};

use roselite_config::{Features, Monitor};
use roselite_request::RoseliteRequest;
use roselite_request::http_caller::HttpCaller;
use roselite_request::icmp_caller::IcmpCaller;

pub fn configure_monitors(monitors: Vec<Monitor>, features: Features) -> Vec<JoinHandle<()>> {
    let mut handles: Vec<JoinHandle<()>> = vec![];

    // Build dependency for http_caller and icmp_caller
    let http_caller = Box::new(HttpCaller::new(features.enable_semyi_status_type));
    let icmp_caller = Box::new(IcmpCaller::new());

    // Start the monitors
    for monitor in monitors {
        println!("Starting monitor for {}", monitor.monitor_target);
        let http_copy = http_caller.clone();
        let icmp_copy = icmp_caller.clone();

        handles.push(spawn(async move {
            loop {
                let tx_ctx = sentry::TransactionContext::new(
                    "Start Roselite monitor",
                    "roselite.start_monitor",
                );
                let transaction = sentry::start_transaction(tx_ctx);

                // Bind the transaction / span to the scope:
                sentry::configure_scope(|scope| scope.set_span(Some(transaction.clone().into())));

                let request = RoseliteRequest::new(http_copy.clone(), icmp_copy.clone());
                let current_time = Instant::now();
                if let Err(err) = request.perform_task(monitor.clone()).await {
                    // Do nothing of this error
                    capture_anyhow(&err);
                    eprintln!("Unexpected error during performing task: {}", err);
                }

                transaction.finish();

                let elapsed = current_time.elapsed();
                if 60 - elapsed.as_secs() > 0 {
                    let sleeping_duration = Duration::from_secs(60 - elapsed.as_secs());
                    println!(
                        "Monitor for {0} will be sleeping for {1} seconds",
                        monitor.monitor_target,
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

    handles
}
