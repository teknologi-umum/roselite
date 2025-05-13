package roselite

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/guregu/null/v6"
	"github.com/teknologi-umum/roselite/internal/sentryhttpclient"
)

type HttpCaller struct {
	Client *http.Client
}

// Call implements Caller.
func (h *HttpCaller) Call(ctx context.Context, monitor Monitor) (Heartbeat, error) {
	span := sentry.StartSpan(ctx, "function", sentry.WithDescription("HttpCaller.Call"))
	ctx = span.Context()
	defer span.Finish()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, monitor.MonitorTarget, nil)
	if err != nil {
		return Heartbeat{}, err
	}

	// Custom user agent. It does not matter if it got overwritten by the user.
	request.Header.Set("User-Agent", "Roselite/1.0 (compatible; +https://github.com/teknologi-umum/roselite)")

	for key, value := range monitor.RequestHeaders {
		request.Header.Set(key, value)
	}

	roundTripper := http.DefaultTransport
	if monitor.TLSConfig != nil {
		roundTripper = &http.Transport{
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
			TLSClientConfig:       monitor.TLSConfig,
		}
	}

	client := &http.Client{Transport: roundTripper}
	if monitor.EnableSentrySampling {
		client.Transport = sentryhttpclient.NewSentryRoundTripper(roundTripper)
	}
	// Use the user-supplied one if it's non nil
	if h.Client != nil {
		client = h.Client
	}

	currentInstant := time.Now()
	response, err := client.Do(request)
	if err != nil {
		elapsed := time.Since(currentInstant)
		var httpProtocol, tlsVersion, tlsCipherName string
		var tlsExpiryDate time.Time
		if response != nil {
			httpProtocol = response.Proto
			if response.TLS != nil {
				tlsVersion = tls.VersionName(response.TLS.Version)
				tlsCipherName = tls.CipherSuiteName(response.TLS.CipherSuite)
				if len(response.TLS.PeerCertificates) > 0 {
					cert := response.TLS.PeerCertificates[0]
					if cert != nil {
						tlsExpiryDate = cert.NotAfter
					}
				}
			}
		}
		return Heartbeat{
			Status:            HeartbeatStatusDown,
			Latency:           int64(elapsed.Seconds()),
			AdditionalMessage: null.StringFrom(err.Error()),
			HttpProtocol:      null.NewString(httpProtocol, httpProtocol != ""),
			TLSVersion:        null.NewString(tlsVersion, tlsVersion != ""),
			TLSCipherName:     null.NewString(tlsCipherName, tlsCipherName != ""),
			TLSExpiryDate:     null.NewTime(tlsExpiryDate, !tlsExpiryDate.IsZero()),
		}, err
	}

	elapsed := time.Since(currentInstant)
	ok := HeartbeatStatusUp
	// everything from 2xx-3xx is considered ok
	if response.StatusCode >= http.StatusBadRequest {
		ok = HeartbeatStatusDown
	}

	var httpProtocol, tlsVersion, tlsCipherName string
	var tlsExpiryDate time.Time
	if response != nil {
		httpProtocol = response.Proto
		if response.TLS != nil {
			tlsVersion = tls.VersionName(response.TLS.Version)
			tlsCipherName = tls.CipherSuiteName(response.TLS.CipherSuite)
			if len(response.TLS.PeerCertificates) > 0 {
				cert := response.TLS.PeerCertificates[0]
				if cert != nil {
					tlsExpiryDate = cert.NotAfter
				}
			}
		}
	}

	return Heartbeat{
		Status:            ok,
		Latency:           int64(elapsed.Seconds()),
		AdditionalMessage: null.String{},
		HttpProtocol:      null.NewString(httpProtocol, httpProtocol != ""),
		TLSVersion:        null.NewString(tlsVersion, tlsVersion != ""),
		TLSCipherName:     null.NewString(tlsCipherName, tlsCipherName != ""),
		TLSExpiryDate:     null.NewTime(tlsExpiryDate, !tlsExpiryDate.IsZero()),
	}, nil
}

var _ Caller = (*HttpCaller)(nil)
