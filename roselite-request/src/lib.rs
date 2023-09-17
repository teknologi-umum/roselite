use std::time::Duration;

use anyhow::Result;
use reqwest::{Client, Method, StatusCode, Url};

use roselite_common::heartbeat::Heartbeat;
use roselite_config::{Monitor, MonitorType};

use crate::bonk_caller::BonkCaller;
use crate::http_caller::HttpCaller;
use crate::icmp_caller::IcmpCaller;

mod bonk_caller;
pub mod http_caller;
pub mod icmp_caller;

pub trait RequestCaller: Send + Sync {
    fn call(&self, monitor: Monitor) -> Result<Heartbeat>;
}

pub struct RoseliteRequest {
    http_caller: Box<dyn RequestCaller>,
    icmp_caller: Box<dyn RequestCaller>,
}

unsafe impl Send for RoseliteRequest {}

unsafe impl Sync for RoseliteRequest {}

impl Default for RoseliteRequest {
    fn default() -> Self {
        RoseliteRequest {
            http_caller: Box::new(BonkCaller::new()),
            icmp_caller: Box::new(BonkCaller::new()),
        }
    }
}

impl RoseliteRequest {
    pub fn new(http_caller: Box<HttpCaller>, icmp_caller: Box<IcmpCaller>) -> Self {
        RoseliteRequest {
            http_caller,
            icmp_caller,
        }
    }

    pub async fn call_kuma_endpoint(
        &self,
        upstream_url: String,
        heartbeat: Heartbeat,
    ) -> Result<()> {
        // Retrieve the currently running span
        let parent_span = sentry::configure_scope(|scope| scope.get_span());

        let span: sentry::TransactionOrSpan = match &parent_span {
            Some(parent) => parent
                .start_child(
                    "request.call_kuma_endpoint",
                    "Call upstream Uptime Kuma endpoint",
                )
                .into(),
            None => {
                let ctx = sentry::TransactionContext::new(
                    "Call upstream Uptime Kuma endpoint",
                    "request.call_kuma_endpoint",
                );
                sentry::start_transaction(ctx).into()
            }
        };

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
                span.clone().finish();
                Err(err)
            }
        }?;

        span.finish();
        Ok(())
    }

    /// It calls the monitor endpoint to create a heartbeat that will be sent to the
    /// push endpoint.
    pub async fn perform_task(&self, monitor: Monitor) -> Result<Heartbeat> {
        // Retrieve the currently running span
        let parent_span = sentry::configure_scope(|scope| scope.get_span());

        let span: sentry::TransactionOrSpan = match &parent_span {
            Some(parent) => parent
                .start_child("request.perform_task", "Perform request task")
                .into(),
            None => {
                let ctx =
                    sentry::TransactionContext::new("Perform request task", "request.perform_task");
                sentry::start_transaction(ctx).into()
            }
        };

        let heartbeat = match monitor.monitor_type {
            MonitorType::HTTP => self.http_caller.call(monitor.clone())?,
            MonitorType::ICMP => self.icmp_caller.call(monitor.clone())?,
        };

        self.call_kuma_endpoint(monitor.clone().push_url, heartbeat.clone())
            .await?;

        span.finish();

        Ok(heartbeat)
    }
}
