use std::sync::Arc;

use axum::Router;
use axum::routing::get;

use crate::config::ServerConfig;
use crate::handlers::ping::ping;
use crate::handlers::remote_write::remote_write;

pub fn register_routes(config: ServerConfig) -> Router {
    let server_config = Arc::new(config);
    Router::new()
        .route("/ping", get(ping))
        .route("/api/push/:id", get(remote_write).post(remote_write))
        .with_state(server_config)
}
