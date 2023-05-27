use crate::handlers::ping::ping;
use axum::routing::get;
use axum::Router;

pub fn register_routes() -> Router {
    Router::new().route("/ping", get(ping))
}
