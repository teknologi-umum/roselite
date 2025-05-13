package main

type ErrorReporting struct {
	SentryDSN              string  `json:"sentry_dsn" toml:"sentry_dsn" yaml:"sentry_dsn" env:"SENTRY_DSN"`
	SentrySampleRate       float64 `json:"sentry_sample_rate" toml:"sentry_sample_rate" yaml:"sentry_sample_rate" env:"SENTRY_SAMPLE_RATE" default:"1.0"`
	SentryTracesSampleRate float64 `json:"sentry_traces_sample_rate" toml:"sentry_traces_sample_rate" yaml:"sentry_traces_sample_rate" env:"SENTRY_TRACES_SAMPLE_RATE" default:"1.0"`
}

type TLSConfig struct {
	CertificateAuthorityFile string `json:"ca_file" toml:"ca_file" yaml:"ca_file"`
	CertificateFile          string `json:"certificate_file" toml:"certificate_file" yaml:"certificate_file"`
	PrivateKeyFile           string `json:"private_key_file" toml:"private_key_file" yaml:"private_key_file"`
	SkipTLSVerify            bool   `json:"skip_tls_verify" toml:"skip_tls_verify" yaml:"skip_tls_verify"`
}

type ServerConfig struct {
	ListenAddress          string            `json:"listen_address" toml:"listen_address" yaml:"listen_address" env:"LISTEN_ADDRESS" default:"127.0.0.1:8321"`
	UpstreamKuma           string            `json:"upstream_kuma" toml:"upstream_kuma" yaml:"upstream_kuma" env:"UPSTREAM_KUMA"`
	UpstreamRequestHeaders map[string]string `json:"upstream_request_headers" toml:"upstream_request_headers" yaml:"upstream_request_headers"`
	UpstreamTLSConfig      TLSConfig         `json:"upstream_tls_config" toml:"upstream_tls_config" yaml:"upstream_tls_config"`
}

type Monitor struct {
	MonitorType    string            `json:"monitor_type" toml:"monitor_type" yaml:"monitor_type"`
	PushURL        string            `json:"push_url" toml:"push_url" yaml:"push_url"`
	MonitorTarget  string            `json:"monitor_target" toml:"monitor_target" yaml:"monitor_target"`
	RequestHeaders map[string]string `json:"request_headers" toml:"request_headers" yaml:"request_headers"`
	// Deprecated: Use TLSConfig instead
	SkipTLSVerify bool      `json:"skip_tls_verify" toml:"skip_tls_verify" yaml:"skip_tls_verify"`
	TLSConfig     TLSConfig `json:"tls_config" toml:"tls_config" yaml:"tls_config"`
}

type Configuration struct {
	ErrorReporting ErrorReporting `json:"error_reporting" toml:"error_reporting" yaml:"error_reporting"`
	ServerConfig   ServerConfig   `json:"server" toml:"server" yaml:"server"`
	Monitors       []Monitor      `json:"monitors" toml:"monitors" yaml:"monitors"`
}
