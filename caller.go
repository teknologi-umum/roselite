package roselite

import "context"

type Caller interface {
	// Call performs a heartbeat call to the monitor endpoint.
	// The implementation of this method should be thread-safe.
	Call(context.Context, Monitor) (Heartbeat, error)
}
