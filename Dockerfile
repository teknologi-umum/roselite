FROM rust:1.85-bookworm AS builder
WORKDIR /app
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y pkg-config libudev-dev perl libssl-dev
COPY . .
RUN cargo build --release

FROM debian:bookworm-slim AS runtime
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/release/roselite /usr/local/bin/roselite
CMD ["/usr/local/bin/roselite"]