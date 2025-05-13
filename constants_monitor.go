package roselite

import (
	"crypto/tls"
	"time"
)

type Monitor struct {
	ID                   string            `json:"id" toml:"id" yaml:"id"`
	MonitorType          MonitorType       `json:"monitor_type" toml:"monitor_type" yaml:"monitor_type"`
	PushURL              string            `json:"push_url" toml:"push_url" yaml:"push_url"`
	MonitorTarget        string            `json:"monitor_target" toml:"monitor_target" yaml:"monitor_target"`
	RequestHeaders       map[string]string `json:"request_headers" toml:"request_headers" yaml:"request_headers"`
	TLSConfig            *tls.Config       `json:"tls_config" toml:"tls_config" yaml:"tls_config"`
	Interval             time.Duration     `json:"interval" toml:"interval" yaml:"interval"`
	EnableSentrySampling bool              `json:"enable_sentry_sampling" toml:"enable_sentry_sampling" yaml:"enable_sentry_sampling"`
}
