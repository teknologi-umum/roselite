package roselite

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/getsentry/sentry-go"
)

func callKumaEndpoint(ctx context.Context, upstreamKumaAddress string, upstreamRequestHeaders map[string]string, httpClient *http.Client, id string, heartbeat Heartbeat) error {
	span := sentry.StartSpan(ctx, "function", sentry.WithDescription("callKumaEndpoint"))
	ctx, cancel := context.WithTimeout(span.Context(), time.Minute*5)
	defer cancel()
	defer span.Finish()

	requestUrl, err := url.JoinPath(upstreamKumaAddress, "/api/push/"+id)
	if err != nil {
		return fmt.Errorf("joining path: %w", err)
	}
	requestUrl = requestUrl + "?" + heartbeat.ToQuery().Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	// Custom user agent. It does not matter if it got overwritten by the user.
	request.Header.Set("User-Agent", "Roselite/1.0 (compatible; +https://github.com/teknologi-umum/roselite)")

	for key, value := range upstreamRequestHeaders {
		request.Header.Set(key, value)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("performing request: %w", err)
	}
	defer func() {
		if response.Body != nil {
			_ = response.Body.Close()
		}
	}()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
