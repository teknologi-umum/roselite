use std::time::Duration;

use anyhow::Result;
use reqwest::{Client, Method, StatusCode, Url};
use tokio::time::Instant;

use roselite_config::Monitor;

/// perform_task calls the monitor endpoint to create a heartbeat that will be sent to the
/// push endpoint.
pub async fn perform_task(monitor: Monitor) -> Result<()> {
    let monitor_client = Client::builder().user_agent("Roselite/1.0").build()?;
    let push_client = Client::builder().user_agent("Roselite/1.0").build()?;

    let current_instant = Instant::now();
    let response = monitor_client
        .request(Method::GET, monitor.monitor_url.as_str())
        .timeout(Duration::from_secs(30))
        .send()
        .await?;

    let elapsed: Duration = current_instant.elapsed();

    let status_code: StatusCode = response.status();
    let mut ok = true;
    // everything from 2xx-3xx is considered ok
    if status_code >= StatusCode::BAD_REQUEST {
        ok = false;
    }

    let mut push_url = Url::parse(monitor.push_url.as_str())?;
    push_url
        .query_pairs_mut()
        .append_pair("msg", "OK")
        .append_pair("status", if ok { "up" } else { "down" })
        .append_pair("ping", elapsed.as_millis().to_string().as_str());

    match push_client
        .request(Method::GET, push_url)
        .timeout(Duration::from_secs(60))
        .send()
        .await
    {
        Ok(response) => {
            if response.status() >= StatusCode::BAD_REQUEST {
                println!(
                    "Received response status of {} during sending event to remote push url",
                    response.status()
                );
                if let Ok(body) = response.text().await {
                    println!("Response body: {}", body)
                }
                return Ok(());
            }
            println!("Successfully sent event to remote push url");
            Ok(())
        }
        Err(err) => {
            println!(
                "An error occurred during sending event to remote push url: {}",
                err
            );
            Err(err)
        }
    }?;

    Ok(())
}
