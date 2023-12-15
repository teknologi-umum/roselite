use std::net::SocketAddr;

use anyhow::Result;
use axum::serve;
use tokio::net::TcpListener;

use crate::config::ServerConfig;
use crate::routes::register_routes;

pub mod config;
mod handlers;
mod routes;

pub async fn run(config: ServerConfig) -> Result<()> {
    // TODO: handle server running with TLS
    let app = register_routes(config.clone());
    let socket_address: SocketAddr = config.address.parse().unwrap();
    let listener = TcpListener::bind(&socket_address).await.unwrap();
    serve(listener, app).await?;

    Ok(())
}
