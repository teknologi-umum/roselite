package roselite_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/guregu/null/v6"
	"github.com/teknologi-umum/roselite"
)

func TestHeartbeatFromQuery(t *testing.T) {
	t.Run("Uptime Kuma only", func(t *testing.T) {
		query := url.Values{}
		query.Set("status", "up")
		query.Set("ping", "100")
		query.Set("msg", "message")

		heartbeat := roselite.HeartbeatFromQuery(query)
		if heartbeat.Status != roselite.HeartbeatStatusUp {
			t.Errorf("expected status to be up, got %s", heartbeat.Status)
		}
		if heartbeat.Latency != 100 {
			t.Errorf("expected latency to be 100, got %d", heartbeat.Latency)
		}
		if heartbeat.AdditionalMessage.ValueOrZero() != "message" {
			t.Errorf("expected additional message to be message, got %s", heartbeat.AdditionalMessage.ValueOrZero())
		}
		if heartbeat.HttpProtocol.Valid {
			t.Errorf("expected http protocol to be invalid, got %s", heartbeat.HttpProtocol.ValueOrZero())
		}
		if heartbeat.TLSVersion.Valid {
			t.Errorf("expected tls version to be invalid, got %s", heartbeat.TLSVersion.ValueOrZero())
		}
		if heartbeat.TLSCipherName.Valid {
			t.Errorf("expected tls cipher name to be invalid, got %s", heartbeat.TLSCipherName.ValueOrZero())
		}
		if heartbeat.TLSExpiryDate.Valid {
			t.Errorf("expected tls expiry date to be invalid, got %s", heartbeat.TLSExpiryDate.ValueOrZero())
		}
	})

	t.Run("status only", func(t *testing.T) {
		query := url.Values{}
		query.Set("status", "up")

		heartbeat := roselite.HeartbeatFromQuery(query)
		if heartbeat.Status != roselite.HeartbeatStatusUp {
			t.Errorf("expected status to be up, got %s", heartbeat.Status)
		}
		if heartbeat.Latency != 0 {
			t.Errorf("expected latency to be 0, got %d", heartbeat.Latency)
		}
		if heartbeat.AdditionalMessage.Valid {
			t.Errorf("expected additional message to be invalid, got %s", heartbeat.AdditionalMessage.ValueOrZero())
		}
	})

	t.Run("Semyi compatible", func(t *testing.T) {
		query := url.Values{}
		query.Set("status", "up")
		query.Set("ping", "100")
		query.Set("msg", "message")
		query.Set("http_protocol", "HTTP/1.1")
		query.Set("tls_version", "TLSv1.3")
		query.Set("tls_cipher", "TLS_AES_256_GCM_SHA384")
		query.Set("tls_expiry", "1677721600")

		heartbeat := roselite.HeartbeatFromQuery(query)
		if heartbeat.Status != roselite.HeartbeatStatusUp {
			t.Errorf("expected status to be up, got %s", heartbeat.Status)
		}
		if heartbeat.Latency != 100 {
			t.Errorf("expected latency to be 100, got %d", heartbeat.Latency)
		}
		if heartbeat.AdditionalMessage.ValueOrZero() != "message" {
			t.Errorf("expected additional message to be message, got %s", heartbeat.AdditionalMessage.ValueOrZero())
		}
		if heartbeat.HttpProtocol.ValueOrZero() != "HTTP/1.1" {
			t.Errorf("expected http protocol to be HTTP/1.1, got %s", heartbeat.HttpProtocol.ValueOrZero())
		}
		if heartbeat.TLSVersion.ValueOrZero() != "TLSv1.3" {
			t.Errorf("expected tls version to be TLSv1.3, got %s", heartbeat.TLSVersion.ValueOrZero())
		}
		if heartbeat.TLSCipherName.ValueOrZero() != "TLS_AES_256_GCM_SHA384" {
			t.Errorf("expected tls cipher name to be TLS_AES_256_GCM_SHA384, got %s", heartbeat.TLSCipherName.ValueOrZero())
		}
		if heartbeat.TLSExpiryDate.Time.Unix() != 1677721600 {
			t.Errorf("expected tls expiry date to be 1677721600, got %d", heartbeat.TLSExpiryDate.Time.Unix())
		}
	})
}

func TestHeartbeat_ToQuery(t *testing.T) {
	t.Run("Uptime Kuma only", func(t *testing.T) {
		heartbeat := roselite.Heartbeat{
			Status:            roselite.HeartbeatStatusUp,
			Latency:           100,
			AdditionalMessage: null.NewString("message", true),
		}

		query := heartbeat.ToQuery()
		if query.Get("status") != "up" {
			t.Errorf("expected status to be up, got %s", query.Get("status"))
		}
		if query.Get("ping") != "100" {
			t.Errorf("expected ping to be 100, got %s", query.Get("ping"))
		}
		if query.Get("msg") != "message" {
			t.Errorf("expected message to be message, got %s", query.Get("msg"))
		}
	})

	t.Run("Semyi compatible", func(t *testing.T) {
		heartbeat := roselite.Heartbeat{
			Status:            roselite.HeartbeatStatusUp,
			Latency:           100,
			AdditionalMessage: null.NewString("message", true),
			HttpProtocol:      null.NewString("HTTP/1.1", true),
			TLSVersion:        null.NewString("TLSv1.3", true),
			TLSCipherName:     null.NewString("TLS_AES_256_GCM_SHA384", true),
			TLSExpiryDate:     null.NewTime(time.Unix(1677721600, 0), true),
		}

		query := heartbeat.ToQuery()
		if query.Get("status") != "up" {
			t.Errorf("expected status to be up, got %s", query.Get("status"))
		}
		if query.Get("ping") != "100" {
			t.Errorf("expected ping to be 100, got %s", query.Get("ping"))
		}
		if query.Get("msg") != "message" {
			t.Errorf("expected message to be message, got %s", query.Get("msg"))
		}
		if query.Get("http_protocol") != "HTTP/1.1" {
			t.Errorf("expected http protocol to be HTTP/1.1, got %s", query.Get("http_protocol"))
		}
		if query.Get("tls_version") != "TLSv1.3" {
			t.Errorf("expected tls version to be TLSv1.3, got %s", query.Get("tls_version"))
		}
		if query.Get("tls_cipher") != "TLS_AES_256_GCM_SHA384" {
			t.Errorf("expected tls cipher name to be TLS_AES_256_GCM_SHA384, got %s", query.Get("tls_cipher"))
		}
		if query.Get("tls_expiry") != "1677721600" {
			t.Errorf("expected tls expiry date to be 1677721600, got %s", query.Get("tls_expiry"))
		}
	})
}
