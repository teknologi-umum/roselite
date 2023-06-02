use crate::config::ServerConfig;
use crate::routes::register_routes;
use anyhow::Result;
use axum::Server;
use std::net::SocketAddr;

pub mod config;
mod handlers;
mod routes;

pub async fn run(config: ServerConfig) -> Result<()> {
    // TODO: handle server running with TLS
    let app = register_routes(config.clone());
    let socket_address: SocketAddr = config.address.parse().unwrap();

    Server::bind(&socket_address)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}
