package roselite

import (
	"context"

	"github.com/guregu/null/v6"
)

type NoopCaller struct {
}

// Call implements Caller.
func (n *NoopCaller) Call(ctx context.Context, monitor Monitor) (Heartbeat, error) {
	return Heartbeat{
		Status:            HeartbeatStatusUp,
		Latency:           0,
		AdditionalMessage: null.String{},
		HttpProtocol:      null.String{},
		TLSVersion:        null.String{},
		TLSCipherName:     null.String{},
		TLSExpiryDate:     null.Time{},
	}, nil
}

var _ Caller = (*NoopCaller)(nil)
