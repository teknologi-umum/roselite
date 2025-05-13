package roselite

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/teknologi-umum/roselite/internal/sentryhttpclient"
)

type Agent struct {
	wg                     *sync.WaitGroup
	shutdownCtx            context.Context
	shutdownCancel         context.CancelFunc
	httpClient             *http.Client
	upstreamRequestHeaders map[string]string
	upstreamKumaAddress    string
}

type AgentOptions struct {
	Monitors               []Monitor
	UpstreamKumaAddress    string
	UpstreamRequestHeaders map[string]string
	UpstreamTLSConfig      *tls.Config
	RegionIdentifier       string
}

var _ io.Closer = (*Agent)(nil)

func NewAgent(options AgentOptions) *Agent {
	httpClientTransport := &http.Transport{
		// Adapted from http.DefaultTransport
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       options.UpstreamTLSConfig,
	}

	httpClient := &http.Client{
		Transport: sentryhttpclient.NewSentryRoundTripper(httpClientTransport),
		Timeout:   time.Minute * 3,
	}

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(sentry.SetHubOnContext(context.Background(), sentry.CurrentHub()))

	upstreamRequestHeaders := make(map[string]string)
	if options.UpstreamRequestHeaders != nil {
		for key, value := range options.UpstreamRequestHeaders {
			upstreamRequestHeaders[key] = value
		}
	}
	region := options.RegionIdentifier
	if region == "" {
		region = "default"
	}

	upstreamRequestHeaders["X-Roselite-Region"] = region
	a := &Agent{
		upstreamKumaAddress:    options.UpstreamKumaAddress,
		upstreamRequestHeaders: upstreamRequestHeaders,
		httpClient:             httpClient,
		wg:                     wg,
		shutdownCtx:            ctx,
		shutdownCancel:         cancel,
	}

	for _, monitor := range options.Monitors {
		wg.Add(1)
		go func(monitor Monitor) {
			defer wg.Done()
			select {
			case <-a.shutdownCtx.Done():
				return
			case <-time.After(monitor.Interval):
				ctx, cancel := context.WithTimeout(a.shutdownCtx, time.Minute*5)
				span := sentry.StartSpan(ctx, "function", sentry.WithDescription("Agent.Start.monitor.loop"))
				span.SetData("roselite.monitor.id", monitor.ID)
				span.SetData("roselite.monitor.type", monitor.MonitorType.String())

				var heartbeat Heartbeat
				var err error
				switch monitor.MonitorType {
				case MonitorTypeHTTP:
					httpCaller := &HttpCaller{}
					heartbeat, err = httpCaller.Call(ctx, monitor)
					break
				case MonitorTypeICMP:
					icmpCaller := &IcmpCaller{Privileged: false}
					heartbeat, err = icmpCaller.Call(ctx, monitor)
					break
				default:
					err = fmt.Errorf("unknown monitor type: %s", monitor.MonitorType.String())
					break
				}

				// Although it may be an error, the Heartbeat struct must not be empty, we must still send it to
				// the upstream instance.
				if err != nil {
					sentry.GetHubFromContext(ctx).CaptureException(err)
				}

				err = callKumaEndpoint(ctx, a.upstreamKumaAddress, a.upstreamRequestHeaders, a.httpClient, monitor.ID, heartbeat)
				if err != nil {
					sentry.GetHubFromContext(ctx).CaptureException(err)
				}
				span.Finish()
				cancel()
			}
		}(monitor)
	}

	return a
}

func (a *Agent) Start() error {
	a.wg.Wait()
	return nil
}

func (a *Agent) Close() error {
	a.shutdownCancel()

	return nil
}
