use std::net::SocketAddr;
use crate::routes::register_routes;
use anyhow::Result;
use axum::Server;
use tokio::net::TcpListener;

mod handlers;
mod routes;

pub async fn run(address: String) -> Result<()> {
    // TODO: handle server running with TLS
    let app = register_routes();
    let socket_address: SocketAddr = address.parse().unwrap();

    Server::bind(&socket_address)
        .serve(app.into_make_service())
        .await?;

    Ok(())
}
