package roselite

import (
    "context"

    "github.com/getsentry/sentry-go"
    "github.com/guregu/null/v6"
    probing "github.com/prometheus-community/pro-bing"
)

type IcmpCaller struct {
    Privileged bool
}

// Call implements Caller.
func (i *IcmpCaller) Call(ctx context.Context, monitor Monitor) (Heartbeat, error) {
    span := sentry.StartSpan(ctx, "function", sentry.WithDescription("IcmpCaller.Call"))
    ctx = span.Context()
    defer span.Finish()

    pinger, err := probing.NewPinger(monitor.MonitorTarget)
    if err != nil {
        return Heartbeat{
            Status:            HeartbeatStatusDown,
            AdditionalMessage: null.StringFrom(err.Error()),
        }, err
    }
    pinger.SetPrivileged(i.Privileged)
    pinger.Count = 3
    err = pinger.RunWithContext(ctx) // Blocks until finished.
    if err != nil {
        return Heartbeat{
            Status:            HeartbeatStatusDown,
            AdditionalMessage: null.StringFrom(err.Error()),
        }, err
    }

    stats := pinger.Statistics()
    if stats == nil {
        return Heartbeat{
            Status:            HeartbeatStatusDown,
            AdditionalMessage: null.StringFrom("no statistic data is available"),
        }, nil

    }
    ok := stats.PacketLoss <= 0.1

    status := HeartbeatStatusUp
    if !ok {
        status = HeartbeatStatusDown
    }
    return Heartbeat{
        Status:            status,
        Latency:           int64(stats.AvgRtt.Seconds()),
        AdditionalMessage: null.String{},
        HttpProtocol:      null.String{},
        TLSVersion:        null.String{},
        TLSCipherName:     null.String{},
        TLSExpiryDate:     null.Time{},
    }, nil
}

var _ Caller = (*IcmpCaller)(nil)
