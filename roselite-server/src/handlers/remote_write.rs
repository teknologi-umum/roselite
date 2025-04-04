use std::sync::Arc;

use anyhow::Result;
use axum::Json;
use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use reqwest::Url;
use sentry::integrations::anyhow::capture_anyhow;
use sentry::{Level, capture_message};
use serde::{Deserialize, Serialize};

use roselite_common::heartbeat::Heartbeat;
use roselite_request::RoseliteRequest;

use crate::config::{ServerConfig, ServerOptions};

#[derive(Clone, Serialize, Deserialize)]
pub struct RemoteWriteResponse {
    pub ok: bool,
}

pub async fn remote_write(
    Path(id): Path<String>,
    Query(params): Query<Heartbeat>,
    State(server_config): State<Arc<ServerConfig>>,
) -> impl IntoResponse {
    let request = Arc::new(RoseliteRequest::default());
    let upstream_kuma = server_config.get_upstream_kuma();
    if let Some(upstream_url) = upstream_kuma {
        match convert_to_upstream(upstream_url, id) {
            Ok(push_url) => match request.call_kuma_endpoint(push_url, params).await {
                Ok(_) => (StatusCode::OK, Json(RemoteWriteResponse { ok: true })),
                Err(error) => {
                    capture_anyhow(&error);

                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Json(RemoteWriteResponse { ok: false }),
                    )
                }
            },
            Err(error) => {
                capture_anyhow(&error);

                (
                    StatusCode::INTERNAL_SERVER_ERROR,
                    Json(RemoteWriteResponse { ok: false }),
                )
            }
        }
    } else {
        capture_message(
            "Remote write was called, yet upstream_url is empty",
            Level::Error,
        );

        (
            StatusCode::PRECONDITION_FAILED,
            Json(RemoteWriteResponse { ok: false }),
        )
    }
}

pub fn convert_to_upstream(upstream_url: String, id: String) -> Result<String> {
    let base = Url::parse(upstream_url.as_str())?;

    let url = base.join(format!("/api/push/{}", id).as_str())?;

    Ok(url.to_string())
}

#[cfg(test)]
mod tests {
    use crate::handlers::remote_write::convert_to_upstream;

    #[test]
    pub fn should_convert_upstream_url() {
        let upstream_url = String::from("https://uptime-kuma.selfhosted.com");
        let id = String::from("tiMmL99KBB");

        let resulting_url = convert_to_upstream(upstream_url, id);

        assert!(resulting_url.is_ok());

        if let Ok(url) = resulting_url {
            assert_eq!(
                url,
                "https://uptime-kuma.selfhosted.com/api/push/tiMmL99KBB"
            );
        }
    }
}
