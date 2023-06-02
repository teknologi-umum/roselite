use std::time::Duration;

use anyhow::Result;
use reqwest::{Client, Method, StatusCode, Url};
use tokio::time::Instant;

use roselite_common::heartbeat::{Heartbeat, HeartbeatStatus};
use roselite_config::Monitor;

pub async fn call_monitor_endpoint(monitor: Monitor) -> Result<Heartbeat> {
    let monitor_client = Client::builder().user_agent("Roselite/1.0").build()?;

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

    Ok(Heartbeat {
        msg: "OK".to_string(),
        status: if ok {
            HeartbeatStatus::Up
        } else {
            HeartbeatStatus::Down
        },
        ping: elapsed.as_millis(),
    })
}

pub async fn call_kuma_endpoint(upstream_url: String, heartbeat: Heartbeat) -> Result<()> {
    let push_client = Client::builder().user_agent("Roselite/1.0").build()?;

    let mut push_url = Url::parse(upstream_url.as_str())?;
    push_url
        .query_pairs_mut()
        .append_pair("msg", heartbeat.msg.as_str())
        .append_pair("status", heartbeat.status.to_string().as_str())
        .append_pair("ping", heartbeat.ping.to_string().as_str());

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
            println!("Successfully sent an event to remote push url");
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

/// It calls the monitor endpoint to create a heartbeat that will be sent to the
/// push endpoint.
pub async fn perform_task(monitor: Monitor) -> Result<Heartbeat> {
    let heartbeat = call_monitor_endpoint(monitor.clone()).await?;

    call_kuma_endpoint(monitor.clone().push_url, heartbeat.clone()).await?;

    Ok(heartbeat)
}
