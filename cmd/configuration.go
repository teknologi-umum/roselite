package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/teknologi-umum/roselite"
)

// ErrorReporting represents the configuration settings for error reporting and monitoring in the application.
type ErrorReporting struct {
	// SentryDSN is the Data Source Name used to configure Sentry for error reporting and monitoring.
	SentryDSN string `json:"sentry_dsn" toml:"sentry_dsn" yaml:"sentry_dsn" env:"SENTRY_DSN"`

	// SentrySampleRate defines the sample rate for sending events to Sentry, represented as a float between 0.0 and 1.0.
	SentrySampleRate float64 `json:"sentry_sample_rate" toml:"sentry_sample_rate" yaml:"sentry_sample_rate" env:"SENTRY_SAMPLE_RATE" default:"1.0"`

	// SentryTracesSampleRate defines the sample rate for tracing events to be sent to Sentry, defaulting to 1.0.
	SentryTracesSampleRate float64 `json:"sentry_traces_sample_rate" toml:"sentry_traces_sample_rate" yaml:"sentry_traces_sample_rate" env:"SENTRY_TRACES_SAMPLE_RATE" default:"1.0"`

	// LogLevel specifies the logging level for the application, with a default value of "info".
	LogLevel slog.Level `json:"log_level" toml:"log_level" yaml:"log_level" env:"LOG_LEVEL" default:"info"`
}

// TLSConfig represents the configuration for TLS, including file paths for certificates and an option to skip verification.
type TLSConfig struct {
	// CertificateAuthorityFile specifies the file path to the certificate authority for verifying server certificates.
	CertificateAuthorityFile string `json:"ca_file" toml:"ca_file" yaml:"ca_file"`

	// CertificateFile specifies the file path to the TLS certificate for secure communication.
	CertificateFile string `json:"certificate_file" toml:"certificate_file" yaml:"certificate_file"`

	// PrivateKeyFile specifies the path to the private key file used for TLS configuration.
	PrivateKeyFile string `json:"private_key_file" toml:"private_key_file" yaml:"private_key_file"`

	// SkipTLSVerify determines whether TLS verification for certificates should be skipped when establishing connections.
	SkipTLSVerify bool `json:"skip_tls_verify" toml:"skip_tls_verify" yaml:"skip_tls_verify"`
}

// ToTLSConfig generates a *tls.Config based on the TLSConfig struct, including certificates and verification settings.
func (t TLSConfig) ToTLSConfig() (*tls.Config, error) {
	var certificates []tls.Certificate = nil
	caCertPool, _ := x509.SystemCertPool()
	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}

	if t.CertificateFile != "" && t.PrivateKeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(t.CertificateFile, t.PrivateKeyFile)
		if err != nil {
			return nil, err
		}

		certificates = make([]tls.Certificate, 0)
		certificates = append(certificates, certificate)
	}

	if t.CertificateAuthorityFile != "" {
		content, err := os.ReadFile(t.CertificateAuthorityFile)
		if err != nil {
			return nil, err
		}
		if ok := caCertPool.AppendCertsFromPEM(content); !ok {
			slog.Warn("failed to append certificate authority")
		}
	}

	return &tls.Config{
		Certificates:       certificates,
		RootCAs:            caCertPool,
		InsecureSkipVerify: t.SkipTLSVerify,
	}, nil
}

// ServerConfig holds the configuration for the server, including listen address, upstream settings, and TLS configuration.
type ServerConfig struct {
	// ListenAddress specifies the address the server listens to on, including IP and port, with a default of 127.0.0.1:8321.
	ListenAddress string `json:"listen_address" toml:"listen_address" yaml:"listen_address" env:"LISTEN_ADDRESS" default:"127.0.0.1:8321"`

	// TLSConfig represents the structure for configuring TLS settings, including certificates and verification options.
	TLSConfig TLSConfig `json:"tls_config" toml:"tls_config" yaml:"tls_config"`

	// UpstreamKuma represents the URL or address of the upstream Kuma service to which requests will be forwarded.
	//
	// Deprecated: Specify UpstreamConfig.BaseUrl instead. The value of this option will be ignored.
	UpstreamKuma string `json:"upstream_kuma" toml:"upstream_kuma" yaml:"upstream_kuma" env:"UPSTREAM_KUMA"`
}

// UpstreamConfig defines the configuration for upstream communication, including base URL, request headers, and TLS settings.
type UpstreamConfig struct {
	// BaseUrl specifies the base URL for upstream requests, supporting JSON, TOML, and YAML configurations.
	BaseUrl string `json:"base_url" toml:"base_url" yaml:"base_url"`

	// RequestHeaders defines a map of headers to be included in upstream requests, with header names as keys and values as values.
	RequestHeaders map[string]string `json:"request_headers" toml:"request_headers" yaml:"request_headers"`

	// TLSConfig represents the structure for configuring TLS settings, including certificates and verification options.
	TLSConfig TLSConfig `json:"tls_config" toml:"tls_config" yaml:"tls_config"`
}

// Monitor represents a monitoring configuration specifying its type, target, interval, request headers, and TLS settings.
type Monitor struct {
	// Id is a unique identifier for the Monitor instance, serialized in JSON, TOML, and YAML formats.
	Id string `json:"id" toml:"id" yaml:"id"`

	// MonitorType represents the type of the monitor, such as HTTP or ICMP, used to define monitoring behavior.
	MonitorType string `json:"monitor_type" toml:"monitor_type" yaml:"monitor_type"`

	// PushURL defines the URL used to send data or updates from the monitor.
	//
	// Deprecated: Specify UpstreamConfig.BaseUrl as the base URL, and Id as the resource path instead.
	PushURL string `json:"push_url" toml:"push_url" yaml:"push_url"`

	// MonitorTarget specifies the target address or resource being monitored.
	MonitorTarget string `json:"monitor_target" toml:"monitor_target" yaml:"monitor_target"`

	// RequestHeaders stores custom headers for HTTP requests as a map of key-value pairs.
	RequestHeaders map[string]string `json:"request_headers" toml:"request_headers" yaml:"request_headers"`

	// SkipTLSVerify indicates whether to bypass TLS certificate verification for secure connections.
	//
	// Deprecated: Use TLSConfig instead. The value of this option will be ignored.
	SkipTLSVerify bool `json:"skip_tls_verify" toml:"skip_tls_verify" yaml:"skip_tls_verify"`

	// TLSConfig represents the TLS-related settings, including certificate paths and skip verification options.
	TLSConfig TLSConfig `json:"tls_config" toml:"tls_config" yaml:"tls_config"`

	// Interval specifies the interval in seconds for how often the monitor performs its checks.
	Interval int `json:"interval" toml:"interval" yaml:"interval"`

	// EnableSentrySampling indicates whether Sentry sampling is enabled for reporting errors or monitoring.
	EnableSentrySampling bool `json:"enable_sentry_sampling" toml:"enable_sentry_sampling" yaml:"enable_sentry_sampling"`
}

// ToRoseliteMonitor converts a Monitor instance to a roselite.Monitor, applying necessary transformations and defaults.
func (m Monitor) ToRoseliteMonitor() roselite.Monitor {
	monitorType, err := roselite.MonitorTypeFromString(m.MonitorType)
	if err != nil {
		slog.Warn(fmt.Sprintf("invalid monitor type: %s", m.MonitorType))
	}

	tlsConfig, err := m.TLSConfig.ToTLSConfig()
	if err != nil {
		slog.Warn(fmt.Sprintf("invalid TLS config: %s", err))
	}

	var interval = time.Duration(m.Interval) * time.Second
	if interval <= 0 {
		interval = time.Second * 30
	}

	return roselite.Monitor{
		ID:                   m.MonitorTarget,
		MonitorType:          monitorType,
		PushURL:              m.PushURL,
		MonitorTarget:        m.MonitorTarget,
		RequestHeaders:       m.RequestHeaders,
		TLSConfig:            tlsConfig,
		Interval:             interval,
		EnableSentrySampling: false,
	}
}

// Configuration represents the root configuration structure containing error reporting, server, and monitors settings.
type Configuration struct {
	ErrorReporting ErrorReporting `json:"error_reporting" toml:"error_reporting" yaml:"error_reporting"`

	// ServerConfig holds the configuration for the server, including listen address, upstream settings, and TLS configuration.
	ServerConfig ServerConfig `json:"server" toml:"server" yaml:"server"`

	// UpstreamConfig defines the configuration for upstream communication, including base URL, request headers, and TLS settings.
	UpstreamConfig UpstreamConfig `json:"upstream" toml:"upstream" yaml:"upstream"`

	// Region is the region identifier for the monitor.
	Region string `json:"region" toml:"region" yaml:"region"`

	// Monitors defines a list of monitoring configurations, specifying individual monitor properties and settings.
	Monitors []Monitor `json:"monitors" toml:"monitors" yaml:"monitors"`
}
