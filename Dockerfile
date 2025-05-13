FROM golang:1.24-bookworm AS builder
WORKDIR /app
ENV DEBIAN_FRONTEND=noninteractive
COPY . .
RUN go build -ldflags="-X 'main.version=$(git describe --tags --always)'" -o roselite .

FROM debian:bookworm-slim AS runtime
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y ca-certificates libcap2-bin && rm -rf /var/lib/apt/lists/*
COPY LICENSE /LICENSE
COPY README.md /README.md
COPY --from=builder /app/roselite /usr/local/bin/roselite
RUN setcap cap_net_raw=+ep /usr/local/bin/roselite
CMD ["/usr/local/bin/roselite"]