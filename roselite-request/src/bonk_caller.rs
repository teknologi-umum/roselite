use anyhow::Result;
use roselite_common::heartbeat::{Heartbeat, HeartbeatStatus};
use roselite_config::Monitor;

use crate::RequestCaller;

#[derive(Clone)]
/// BonKCaller is a mock or empty struct that's implement
/// RequestCaller. Do not use this as a normal caller transport
/// unless you know what you're doing.
pub struct BonkCaller {}

impl BonkCaller {
    pub fn new() -> Self {
        BonkCaller {}
    }
}

impl Default for BonkCaller {
    fn default() -> Self {
        Self::new()
    }
}

impl RequestCaller for BonkCaller {
    fn call(&self, _monitor: Monitor) -> Result<Heartbeat> {
        Ok(Heartbeat {
            msg: "OK".to_string(),
            status: HeartbeatStatus::Up,
            ping: 0,
        })
    }
}
