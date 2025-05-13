package roselite

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/teknologi-umum/roselite/internal/sentryhttpclient"
)

type Server struct {
	httpServer             *http.Server
	upstreamKumaAddress    string
	upstreamRequestHeaders map[string]string
	httpClient             *http.Client
}

type ServerOptions struct {
	ListeningAddress       string
	UpstreamKumaAddress    string
	UpstreamRequestHeaders map[string]string
	UpstreamTLSConfig      *tls.Config
	ServerTLSConfig        *tls.Config
}

type remoteWriteResponse struct {
	Ok bool `json:"ok"`
}

func NewServer(options ServerOptions) *Server {
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{})

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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/api/push/{id}", func(w http.ResponseWriter, r *http.Request) {
		if options.UpstreamKumaAddress == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPreconditionFailed)
			_ = json.NewEncoder(w).Encode(remoteWriteResponse{Ok: false})
			return
		}

		id := r.PathValue("id")
		if id == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPreconditionFailed)
			_ = json.NewEncoder(w).Encode(remoteWriteResponse{Ok: false})
			return
		}

		err := callKumaEndpoint(context.WithoutCancel(r.Context()),
			options.UpstreamKumaAddress,
			options.UpstreamRequestHeaders,
			httpClient,
			id,
			HeartbeatFromQuery(r.URL.Query()),
		)
		if err != nil {
			sentry.GetHubFromContext(r.Context()).CaptureException(err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(remoteWriteResponse{Ok: false})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(remoteWriteResponse{Ok: true})
	})
	s := &Server{
		httpServer: &http.Server{
			Addr:              options.ListeningAddress,
			Handler:           sentryMiddleware.Handle(mux),
			TLSConfig:         options.ServerTLSConfig,
			ReadTimeout:       time.Minute,
			ReadHeaderTimeout: time.Minute,
			WriteTimeout:      time.Minute,
			IdleTimeout:       time.Minute,
		},
		httpClient:             httpClient,
		upstreamKumaAddress:    options.UpstreamKumaAddress,
		upstreamRequestHeaders: options.UpstreamRequestHeaders,
	}

	return s
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) ListenAndServeTLS() error {
	s.httpServer.Protocols = new(http.Protocols)
	s.httpServer.Protocols.SetHTTP1(true)
	s.httpServer.Protocols.SetHTTP2(true)
	return s.httpServer.ListenAndServeTLS("", "")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
