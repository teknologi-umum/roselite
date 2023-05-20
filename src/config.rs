use std::collections::BTreeMap;
use std::fs::File;
use std::io::Read;
use serde::Deserialize;
use anyhow::{Error, Result};

#[derive(Deserialize, Clone)]
pub struct Monitor {
    pub push_url: String,
    pub monitor_url: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub request_headers: Option<BTreeMap<String, String>>,
}

#[derive(Deserialize, Clone)]
pub struct Configuration {
    pub monitors: Vec<Monitor>,
}

impl Configuration {
    pub fn from_yaml(value: &str) -> Self {
        serde_yaml::from_str::<Self>(value).unwrap()
    }

    pub fn from_json(value: &str) -> Self {
        json5::from_str::<Self>(value).unwrap()
    }

    pub fn from_toml(value: &str) -> Self {
        toml::from_str::<Self>(value).unwrap()
    }

    pub fn from_file(path: &String) -> Result<Configuration> {
        match File::open(path) {
            Ok(mut file) => {
                let mut contents = String::new();
                file.read_to_string(&mut contents).unwrap();

                if path.ends_with("toml") {
                    return Ok(Self::from_toml(&contents));
                }

                if path.ends_with("json") || path.ends_with("json5") {
                    return Ok(Self::from_json(&contents));
                }

                if path.ends_with("yaml") || path.ends_with("yml") {
                    return Ok(Self::from_yaml(&contents));
                }

                Err(Error::msg("Invalid file type"))
            },
            Err(_) => Err(Error::msg("Failed opening configuration file"))
        }
    }
}