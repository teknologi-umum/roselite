use std::fmt::Display;

use serde::{Deserialize, Serialize, Serializer};

#[derive(Debug, Clone, Deserialize)]
pub enum HeartbeatStatus {
    Up,
    Down,
    DegradedPerformance,
    UnderMaintenance,
    LimitedAvailability,
}

impl Serialize for HeartbeatStatus {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(match self {
            HeartbeatStatus::Up => "up",
            HeartbeatStatus::Down => "down",
            HeartbeatStatus::DegradedPerformance => "degraded_performance",
            HeartbeatStatus::UnderMaintenance => "under_maintenance",
            HeartbeatStatus::LimitedAvailability => "limited_availability",
        })
    }
}

impl Display for HeartbeatStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "{}",
            match self {
                HeartbeatStatus::Up => String::from("up"),
                HeartbeatStatus::Down => String::from("down"),
                HeartbeatStatus::DegradedPerformance => String::from("degraded_performance"),
                HeartbeatStatus::UnderMaintenance => String::from("under_maintenance"),
                HeartbeatStatus::LimitedAvailability => String::from("limited_availability"),
            }
        )
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Heartbeat {
    pub msg: String,
    pub status: HeartbeatStatus,
    pub ping: u128,
}
