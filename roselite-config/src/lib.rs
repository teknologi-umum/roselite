use std::collections::BTreeMap;
use std::fs::File;
use std::io::Read;

use anyhow::{Error, Result};
use serde::Deserialize;

/// Monitor defines a single monitor configuration
#[derive(Deserialize, Clone)]
pub struct Monitor {
    pub push_url: String,
    pub monitor_url: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub request_headers: Option<BTreeMap<String, String>>,
}

/// Configuration sets a global configuration for the application.
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

    /// Supported file extensions (from the given path):
    /// - json, json5
    /// - toml
    /// - yaml, yml
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
            }
            Err(_) => Err(Error::msg("Failed opening configuration file")),
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::Configuration;

    #[test]
    fn parse_toml_configuration() {
        let configuration = r#"[[monitors]]
push_url = "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
monitor_url = "https://github.com/healthz""#;

        let parsed_configuration = Configuration::from_toml(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_url, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }
    }

    #[test]
    fn parse_yaml_configuration() {
        let configuration = r#"monitors:
  - push_url: "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
    monitor_url: "https://github.com/healthz""#;

        let parsed_configuration = Configuration::from_yaml(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_url, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }
    }

    #[test]
    fn parse_json_configuration() {
        let configuration = r#"{
  "monitors": [
    {
        "push_url": "https://your-uptime-kuma.com/api/push/Eq15E23yc3",
        "monitor_url": "https://github.com/healthz"
    }
  ]
}"#;

        let parsed_configuration = Configuration::from_json(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_url, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }
    }
}
