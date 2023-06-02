use roselite_config::Monitor;
use roselite_request::perform_task;
use std::time::Duration;
use tokio::spawn;
use tokio::task::JoinHandle;
use tokio::time::{sleep, Instant};

pub fn configure_monitors(monitors: Vec<Monitor>) -> Vec<JoinHandle<()>> {
    let mut handles: Vec<JoinHandle<()>> = vec![];

    for monitor in monitors {
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

    handles
}
