use std::sync::Arc;

#[derive(Clone)]
pub struct ServerConfig {
    pub address: String,
    pub upstream_kuma: Option<String>,
}

impl ServerOptions for ServerConfig {
    fn get_upstream_kuma(&self) -> Option<String> {
        self.upstream_kuma.clone()
    }
}

pub trait ServerOptions {
    fn get_upstream_kuma(&self) -> Option<String>;
}

pub type DynServerConfig = Arc<dyn ServerOptions + Send + Sync + 'static>;
