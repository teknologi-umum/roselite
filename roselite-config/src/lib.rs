use std::collections::BTreeMap;
use std::fs::File;
use std::io::Read;
use std::str::FromStr;

use anyhow::{Error, Result};
use serde::Deserialize;

#[derive(Debug, Deserialize, Clone, PartialEq)]
pub enum MonitorType {
    HTTP,
    ICMP,
}

impl FromStr for MonitorType {
    type Err = ();

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s {
            "HTTP" => Ok(MonitorType::HTTP),
            "ICMP" => Ok(MonitorType::ICMP),
            _ => Err(()),
        }
    }
}

/// Monitor defines a single monitor configuration
#[derive(Deserialize, Clone, Debug)]
pub struct Monitor {
    pub monitor_type: MonitorType,
    pub push_url: String,
    pub monitor_target: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub request_headers: Option<BTreeMap<String, String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub skip_tls_verify: Option<bool>,
}

#[derive(Deserialize, Clone, Debug)]
pub struct ErrorReporting {
    pub sentry_dsn: String,
}

#[derive(Deserialize, Clone, Debug)]
pub struct ServerConfig {
    pub listen_address: String,
    pub upstream_kuma: Option<String>,
}

/// Configuration sets a global configuration for the application.
#[derive(Deserialize, Clone, Debug)]
pub struct Configuration {
    pub error_reporting: Option<ErrorReporting>,
    pub server: Option<ServerConfig>,
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
    use std::env::temp_dir;
    use std::fs::File;
    use std::io::Write;

    use anyhow::Result;

    use crate::{Configuration, MonitorType};

    #[test]
    fn parse_toml_configuration() {
        let configuration = r#"[[monitors]]
monitor_type = "HTTP"
push_url = "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
monitor_target = "https://github.com/healthz""#;

        let parsed_configuration = Configuration::from_toml(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_type, MonitorType::HTTP);
            assert_eq!(first_monitor.monitor_target, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }
    }

    #[test]
    fn parse_yaml_configuration() {
        let configuration = r#"monitors:
  - monitor_type: "HTTP"
    push_url: "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
    monitor_target: "https://github.com/healthz""#;

        let parsed_configuration = Configuration::from_yaml(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_type, MonitorType::HTTP);
            assert_eq!(first_monitor.monitor_target, "https://github.com/healthz");
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
        "monitor_type": "HTTP",
        "push_url": "https://your-uptime-kuma.com/api/push/Eq15E23yc3",
        "monitor_target": "https://github.com/healthz"
    }
  ]
}"#;

        let parsed_configuration = Configuration::from_json(configuration);

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_type, MonitorType::HTTP);
            assert_eq!(first_monitor.monitor_target, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }
    }

    #[test]
    fn parse_file() -> Result<()> {
        let temporary_directory = temp_dir();
        let temporary_file_path = temporary_directory.join("roselite-config.toml");
        let clone_temp_path = temporary_file_path.clone();
        let mut file = File::create(temporary_file_path)?;
        let configuration = r#"[error_reporting]
sentry_dsn = "https://sentry.io"

[server]
listen_address = "127.0.0.1:8321"

[[monitors]]
monitor_type = "ICMP"
push_url = "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
monitor_target = "https://github.com/healthz"
"#;
        let _ = file.write(configuration.as_bytes())?;
        file.flush()?;

        let path_str = match clone_temp_path.to_str() {
            Some(value) => String::from(value),
            None => String::new(),
        };
        let parsed_configuration = Configuration::from_file(&path_str)?;

        assert!(parsed_configuration.server.is_some());
        if let Some(server) = parsed_configuration.server {
            assert_eq!(server.listen_address, "127.0.0.1:8321");
        }

        assert!(parsed_configuration.error_reporting.is_some());
        if let Some(error_reporting) = parsed_configuration.error_reporting {
            assert_eq!(error_reporting.sentry_dsn, "https://sentry.io");
        }

        assert_eq!(parsed_configuration.monitors.len(), 1);
        if let Some(first_monitor) = parsed_configuration.monitors.first() {
            assert_eq!(first_monitor.monitor_type, MonitorType::ICMP);
            assert_eq!(first_monitor.monitor_target, "https://github.com/healthz");
            assert_eq!(
                first_monitor.push_url,
                "https://your-uptime-kuma.com/api/push/Eq15E23yc3"
            );
        }

        Ok(())
    }
}
