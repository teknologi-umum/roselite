package roselite_test

import (
	"crypto/tls"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/teknologi-umum/roselite"
)

func KumaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
}

func KumaSecureServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
}

func TestServer_Unencrypted(t *testing.T) {
	kumaServer := KumaServer()
	t.Cleanup(func() {
		kumaServer.Close()
	})
	randomPort := 10_000 + rand.IntN(50_000)

	server := roselite.NewServer(roselite.ServerOptions{
		ListeningAddress:       ":" + strconv.FormatInt(int64(randomPort), 10),
		UpstreamKumaAddress:    kumaServer.URL,
		UpstreamRequestHeaders: map[string]string{"X-Forwarded-Proto": "http"},
		UpstreamTLSConfig:      kumaServer.TLS,
		ServerTLSConfig:        nil,
	})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start server: %v", err)
		}
	}()
	t.Cleanup(func() {
		err := server.Shutdown(t.Context())
		if err != nil {
			t.Logf("failed to shutdown server: %v", err)
		}
	})

	serverAddress := "http://localhost:" + strconv.FormatInt(int64(randomPort), 10)

	t.Run("Ping endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/ping", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Read the body, expect "OK" plain text response.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if string(body) != "OK" {
			t.Errorf("unexpected response body: %s", string(body))
			return
		}
	})

	t.Run("Push endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/api/push/12?status=up&ping=0", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Expecting response of `{"ok":true}` with `Content-Type: application/json` header.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if strings.TrimSpace(string(body)) != `{"ok":true}` {
			t.Errorf("unexpected response body: %s", string(body))
		}

		// Expecting `Content-Type: application/json` header.
		if response.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", response.Header.Get("Content-Type"))
		}
	})
}

func TestServer_Encrypted(t *testing.T) {
	kumaServer := KumaSecureServer()
	t.Cleanup(func() {
		kumaServer.Close()
	})
	randomPort := 10_000 + rand.IntN(50_000)

	cert, key := GenerateCert()
	certificate, err := tls.X509KeyPair(cert.Bytes(), key.Bytes())
	if err != nil {
		t.Errorf("failed to generate certificate: %v", err)
		return
	}

	serverTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
	}

	server := roselite.NewServer(roselite.ServerOptions{
		ListeningAddress:       ":" + strconv.FormatInt(int64(randomPort), 10),
		UpstreamKumaAddress:    kumaServer.URL,
		UpstreamRequestHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		UpstreamTLSConfig:      kumaServer.Client().Transport.(*http.Transport).TLSClientConfig,
		ServerTLSConfig:        serverTLSConfig,
	})
	go func() {
		err := server.ListenAndServeTLS()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start server: %v", err)
		}
	}()
	t.Cleanup(func() {
		err := server.Shutdown(t.Context())
		if err != nil {
			t.Logf("failed to shutdown server: %v", err)
		}
	})

	serverAddress := "https://localhost:" + strconv.FormatInt(int64(randomPort), 10)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Minute * 3,
	}

	t.Run("Ping endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/ping", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := client.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Read the body, expect "OK" plain text response.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if string(body) != "OK" {
			t.Errorf("unexpected response body: %s", string(body))
			return
		}
	})

	t.Run("Push endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/api/push/12?status=up&ping=0", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := client.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Expecting a response of `{"ok":true}` with `Content-Type: application/json` header.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if strings.TrimSpace(string(body)) != `{"ok":true}` {
			t.Errorf("unexpected response body: %s", string(body))
		}

		// Expecting `Content-Type: application/json` header.
		if response.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", response.Header.Get("Content-Type"))
		}
	})
}

func TestServer_UpstreamKumaAddressNotSet(t *testing.T) {
	randomPort := 10_000 + rand.IntN(50_000)

	server := roselite.NewServer(roselite.ServerOptions{
		ListeningAddress:       ":" + strconv.FormatInt(int64(randomPort), 10),
		UpstreamKumaAddress:    "",
		UpstreamRequestHeaders: nil,
		UpstreamTLSConfig:      nil,
		ServerTLSConfig:        nil,
	})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start server: %v", err)
		}
	}()

	t.Cleanup(func() {
		err := server.Shutdown(t.Context())
		if err != nil {
			t.Logf("failed to shutdown server: %v", err)
		}
	})

	serverAddress := "http://localhost:" + strconv.FormatInt(int64(randomPort), 10)

	t.Run("Push endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/api/push/12?status=up&ping=0", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusPreconditionFailed {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Expecting a response of `{"ok":false}` with `Content-Type: application/json` header.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if strings.TrimSpace(string(body)) != `{"ok":false}` {
			t.Errorf("unexpected response body: %s", string(body))
		}

		// Expecting `Content-Type: application/json` header.
		if response.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", response.Header.Get("Content-Type"))
		}
	})
}

func TestServer_InvalidUpstreamCertificateAuthority(t *testing.T) {
	kumaServer := KumaSecureServer()
	t.Cleanup(func() {
		kumaServer.Close()
	})
	randomPort := 10_000 + rand.IntN(50_000)

	server := roselite.NewServer(roselite.ServerOptions{
		ListeningAddress:       ":" + strconv.FormatInt(int64(randomPort), 10),
		UpstreamKumaAddress:    kumaServer.URL,
		UpstreamRequestHeaders: nil,
		UpstreamTLSConfig:      nil,
		ServerTLSConfig:        nil,
	})
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start server: %v", err)
		}
	}()
	t.Cleanup(func() {
		err := server.Shutdown(t.Context())
		if err != nil {
			t.Logf("failed to shutdown server: %v", err)
		}
	})

	serverAddress := "http://localhost:" + strconv.FormatInt(int64(randomPort), 10)

	t.Run("Push endpoint", func(t *testing.T) {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, serverAddress+"/api/push/12?status=up&ping=0", nil)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
			return
		}

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Errorf("failed to perform request: %v", err)
			return
		}
		defer func() {
			if response.Body != nil {
				_ = response.Body.Close()
			}
		}()

		if response.StatusCode != http.StatusInternalServerError {
			t.Errorf("unexpected status code: %d", response.StatusCode)
			return
		}

		// Expecting a response of `{"ok":false}` with `Content-Type: application/json` header.
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("failed to read response body: %v", err)
			return
		}

		if strings.TrimSpace(string(body)) != `{"ok":false}` {
			t.Errorf("unexpected response body: %s", string(body))
		}

		// Expecting `Content-Type: application/json` header.
		if response.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", response.Header.Get("Content-Type"))
		}
	})
}
