package roselite_test

import (
	"testing"

	"github.com/teknologi-umum/roselite"
)

func TestNoopCaller(t *testing.T) {
	caller := roselite.NoopCaller{}

	heartbeat, err := caller.Call(nil, roselite.Monitor{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if heartbeat.Status != roselite.HeartbeatStatusUp {
		t.Errorf("expected status to be up, got %s", heartbeat.Status)
	}
}
