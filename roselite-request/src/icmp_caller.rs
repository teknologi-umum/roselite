use std::time::Duration;

use anyhow::{Error, Result};
use fastping_rs::PingResult::{Idle, Receive};
use fastping_rs::Pinger;

use roselite_common::heartbeat::{Heartbeat, HeartbeatStatus};
use roselite_config::Monitor;

use crate::RequestCaller;

#[derive(Debug, Clone)]
pub struct IcmpCaller {}

impl IcmpCaller {
    pub fn new() -> Self {
        return IcmpCaller {};
    }
}

impl RequestCaller for IcmpCaller {
    fn call(&self, monitor: Monitor) -> Result<Heartbeat> {
        // Retrieve the currently running span
        let parent_span = sentry::configure_scope(|scope| scope.get_span());

        let span: sentry::TransactionOrSpan = match &parent_span {
            Some(parent) => parent
                .start_child("icmp_caller.call", "Call target ICMP request")
                .into(),
            None => {
                let ctx =
                    sentry::TransactionContext::new("Call target ICMP request", "icmp_caller.call");
                sentry::start_transaction(ctx).into()
            }
        };

        let (pinger, results) = Pinger::new(None, None).unwrap();

        pinger.add_ipaddr(monitor.monitor_target.as_str());
        pinger.run_pinger();

        let mut ok = true;
        let mut ping_duration = Duration::from_secs(0);
        for _ in 0..4 {
            match results.recv() {
                Ok(result) => match result {
                    Idle { addr: _addr } => {
                        ok = false;
                    }
                    Receive { addr: _addr, rtt } => {
                        ok = true;
                        ping_duration = rtt;
                    }
                },
                Err(err) => {
                    span.finish();
                    return Err(Error::from(err));
                }
            }
        }

        pinger.stop_pinger();

        span.finish();

        Ok(Heartbeat {
            msg: "OK".to_string(),
            status: if ok {
                HeartbeatStatus::Up
            } else {
                HeartbeatStatus::Down
            },
            ping: ping_duration.as_millis(),
        })
    }
}
