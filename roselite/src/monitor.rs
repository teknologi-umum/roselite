use roselite_config::Monitor;
use roselite_request::http_caller::HttpCaller;
use roselite_request::icmp_caller::IcmpCaller;
use roselite_request::{RoseliteRequest};
use std::time::Duration;
use tokio::spawn;
use tokio::task::JoinHandle;
use tokio::time::{sleep, Instant};

pub fn configure_monitors(monitors: Vec<Monitor>) -> Vec<JoinHandle<()>> {
    let mut handles: Vec<JoinHandle<()>> = vec![];

    // Build dependency for http_caller and icmp_caller
    let http_caller = Box::new(HttpCaller::new());
    let icmp_caller = Box::new(IcmpCaller::new());
    let request = RoseliteRequest::new(http_caller, icmp_caller);

    // Start the monitors
    for monitor in monitors {
        println!("Starting monitor for {}", monitor.monitor_target);

        handles.push(spawn(async move {
            let cloned_monitor: Monitor = monitor.clone();

            loop {
                let tx_ctx = sentry::TransactionContext::new(
                    "Start Roselite monitor",
                    "roselite.start_monitor",
                );
                let transaction = sentry::start_transaction(tx_ctx);

                // Bind the transaction / span to the scope:
                sentry::configure_scope(|scope| scope.set_span(Some(transaction.into())));

                let current_time = Instant::now();
                if let Err(err) = request.perform_task(monitor.clone()).await {
                    // Do nothing of this error
                    sentry::capture_error(&err);
                    eprintln!("Unexpected error during performing task: {}", err);
                }

                transaction.finish();

                let elapsed = current_time.elapsed();
                if 60 - elapsed.as_secs() > 0 {
                    let sleeping_duration = Duration::from_secs(60 - elapsed.as_secs());
                    println!(
                        "Monitor for {0} will be sleeping for {1} seconds",
                        cloned_monitor.monitor_target,
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
