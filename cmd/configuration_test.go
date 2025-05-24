package main_test

import (
	"encoding/json"
	"testing"

	main "github.com/teknologi-umum/roselite/cmd"
)

func TestConfiguration(t *testing.T) {
	jsonConfiguration := `{
    "error_reporting": {
        "sentry_dsn": "https://00000000000000@ingest.sentry.io/0",
        "sentry_sample_rate": 1.0,
        "sentry_traces_sample_rate": 1.0
    },
    "server": {
        "listen_address": "127.0.0.1:8123",
        "tls_config": {
            "ca_file": "/etc/ssl/certs/ca-certificates.crt",
            "certificate_file": "/etc/ssl/certs/roselite.crt",
            "private_key_file": "/etc/ssl/certs/roselite.key",
            "skip_tls_verify": false
        }
    },
    "upstream": {
        "base_url": "https://kuma.io",
        "request_headers": {
            "Authorization": "Bearer <PASSWORD>"
        },
        "tls_config": {
            "ca_file": "/etc/ssl/certs/ca-certificates.crt",
            "certificate_file": "/etc/ssl/certs/roselite.crt",
            "private_key_file": "/etc/ssl/certs/roselite.key",
            "skip_tls_verify": true
        }
    },
    "region": "ap-southeast-1",
    "monitors": [
        {
            "id": "1",
            "monitor_type": "HTTP",
            "monitor_target": "https://blog.teknologiumum.com",
            "request_headers": {},
            "tls_config": {},
            "interval": 10,
            "enable_sentry_sampling": true
        }
    ]
}`

	var configuration main.Configuration
	err := json.Unmarshal([]byte(jsonConfiguration), &configuration)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}
