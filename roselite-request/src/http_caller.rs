use std::time::Duration;

use anyhow::Result;
use reqwest::blocking::Client;
use reqwest::{Method, StatusCode};
use tokio::time::Instant;

use roselite_common::heartbeat::{Heartbeat, HeartbeatStatus};
use roselite_config::Monitor;

use crate::RequestCaller;

#[derive(Clone)]
pub struct HttpCaller {
    client: Client,
    pub allow_semyi_status_type: bool,
}

impl HttpCaller {
    pub fn new(allow_semyi_status_type: bool) -> Self {
        HttpCaller {
            client: Client::builder()
                .user_agent("Roselite/1.0")
                .build()
                .unwrap(),
            allow_semyi_status_type: allow_semyi_status_type,
        }
    }
}

impl Default for HttpCaller {
    fn default() -> Self {
        Self::new(false)
    }
}

impl RequestCaller for HttpCaller {
    fn call(&self, monitor: Monitor) -> Result<Heartbeat> {
        // Retrieve the currently running span
        let parent_span = sentry::configure_scope(|scope| scope.get_span());

        let span: sentry::TransactionOrSpan = match &parent_span {
            Some(parent) => parent
                .start_child("http_caller.call", "Call target HTTP request")
                .into(),
            None => {
                let ctx =
                    sentry::TransactionContext::new("Call target HTTP request", "http_caller.call");
                sentry::start_transaction(ctx).into()
            }
        };

        let current_instant = Instant::now();
        let response = self
            .client
            .request(Method::GET, monitor.monitor_target.as_str())
            .timeout(Duration::from_secs(30))
            .send()?;

        let elapsed: Duration = current_instant.elapsed();

        let status_code: StatusCode = response.status();

        let mut status: HeartbeatStatus;
        if self.allow_semyi_status_type {
            if status_code >= StatusCode::OK && status_code <= StatusCode::PERMANENT_REDIRECT {
                status = HeartbeatStatus::Up;
            } else if status_code == StatusCode::FORBIDDEN || status_code == StatusCode::GONE {
                status = HeartbeatStatus::LimitedAvailability;
            } else {
                status = HeartbeatStatus::Down;
            }

            if elapsed.as_millis() > 30_000 {
                status = HeartbeatStatus::DegradedPerformance;
            }
        } else {
            if status_code <= StatusCode::BAD_REQUEST {
                status = HeartbeatStatus::Up;
            } else {
                status = HeartbeatStatus::Down;
            }
        }

        span.finish();

        Ok(Heartbeat {
            msg: "OK".to_string(),
            status: status,
            ping: elapsed.as_millis(),
        })
    }
}
