use serde::{Deserialize, Serialize, Serializer};

#[derive(Debug, Clone, Deserialize)]
pub enum HeartbeatStatus {
    Up,
    Down,
}

impl Serialize for HeartbeatStatus {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
        where
            S: Serializer,
    {
        serializer.serialize_str(match self {
            HeartbeatStatus::Up => "up",
            HeartbeatStatus::Down => "down",
        })
    }
}

impl ToString for HeartbeatStatus {
    fn to_string(&self) -> String {
        match self {
            HeartbeatStatus::Up => String::from("up"),
            HeartbeatStatus::Down => String::from("down"),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Heartbeat {
    pub msg: String,
    pub status: HeartbeatStatus,
    pub ping: u128,
}
