use crate::config::{DynServerConfig, ServerConfig};
use crate::handlers::ping::ping;
use crate::handlers::remote_write::remote_write;
use axum::routing::{get, post};
use axum::Router;
use std::sync::Arc;

pub fn register_routes(config: ServerConfig) -> Router {
    let server_config = Arc::new(config) as DynServerConfig;
    Router::new()
        .route("/ping", get(ping))
        .route("/api/push/:id", post(remote_write))
        .with_state(server_config)
}
