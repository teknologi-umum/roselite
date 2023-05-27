use clap::Command;

pub fn cli() -> Command {
    Command::new("roselite")
        .about("Active relay for Uptime Kuma's push monitor type")
        .subcommand_required(false)
        .allow_external_subcommands(true)
        .arg_required_else_help(true)
        .subcommand(
            Command::new("agent")
                .about("Start Roselite in agent mode, it will not expose any HTTP port")
        )
        .subcommand(
            Command::new("server")
                .about("Start Roselite in server mode, it will expose HTTP port that can be used to be pinged by other Roselite instances")
        )
}
