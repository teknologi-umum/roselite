package roselite_test

import (
	"strings"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/teknologi-umum/roselite"
)

func TestHttpCaller(t *testing.T) {
	ctx := sentry.SetHubOnContext(t.Context(), sentry.CurrentHub().Clone())
	caller := roselite.HttpCaller{}

	t.Run("Example", func(t *testing.T) {
		monitor := roselite.Monitor{
			ID:                   "example",
			MonitorType:          roselite.MonitorTypeHTTP,
			PushURL:              "",
			MonitorTarget:        "https://example.com",
			RequestHeaders:       nil,
			TLSConfig:            nil,
			Interval:             time.Minute,
			EnableSentrySampling: true,
		}

		heartbeat, err := caller.Call(ctx, monitor)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if heartbeat.Status != roselite.HeartbeatStatusUp {
			t.Errorf("expected status to be up, got %s", heartbeat.Status)
		}

		if heartbeat.Latency < 0 {
			t.Errorf("expected latency to be positive, got %d", heartbeat.Latency)
		}
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		monitor := roselite.Monitor{
			ID:                   "invalid.domain",
			MonitorType:          roselite.MonitorTypeHTTP,
			PushURL:              "",
			MonitorTarget:        "https://invalid.domain",
			RequestHeaders:       nil,
			TLSConfig:            nil,
			Interval:             time.Minute,
			EnableSentrySampling: false,
		}

		heartbeat, err := caller.Call(ctx, monitor)
		if err == nil {
			t.Errorf("expected error, got nil")
		} else {
			if !strings.Contains(err.Error(), "no such host") {
				t.Errorf("expected error to contain 'no such host', got %s", err)
			}

			if heartbeat.AdditionalMessage.ValueOrZero() != err.Error() {
				t.Errorf("expected additional message to be %s, got %s", err.Error(), heartbeat.AdditionalMessage.ValueOrZero())
			}
		}

		if heartbeat.Status != roselite.HeartbeatStatusDown {
			t.Errorf("expected status to be down, got %s", heartbeat.Status)
		}

		if heartbeat.Latency < 0 {
			t.Errorf("expected latency to be positive, got %d", heartbeat.Latency)
		}
	})
}
