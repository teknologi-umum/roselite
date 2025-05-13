package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/jinzhu/configor"
	"github.com/teknologi-umum/roselite"
	"github.com/urfave/cli/v3"
)

func ServerAction(ctx context.Context, c *cli.Command) error {
	var configuration Configuration
	err := configor.New(&configor.Config{}).Load(&configuration, c.String("config"))
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	monitors := make([]roselite.Monitor, len(configuration.Monitors))
	for i, monitor := range configuration.Monitors {
		monitors[i] = monitor.ToRoseliteMonitor()
	}

	upstreamTLSConfig, err := configuration.UpstreamConfig.TLSConfig.ToTLSConfig()
	if err != nil {
		return fmt.Errorf("creating TLS config: %w", err)
	}

	serverTLSConfig, err := configuration.ServerConfig.TLSConfig.ToTLSConfig()
	if err != nil {
		return fmt.Errorf("creating TLS config: %w", err)
	}

	server := roselite.NewServer(roselite.ServerOptions{
		ListeningAddress:       configuration.ServerConfig.ListenAddress,
		UpstreamKumaAddress:    configuration.UpstreamConfig.BaseUrl,
		UpstreamRequestHeaders: configuration.UpstreamConfig.RequestHeaders,
		UpstreamTLSConfig:      upstreamTLSConfig,
		ServerTLSConfig:        serverTLSConfig,
	})

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal
		err := server.Shutdown(ctx)
		if err != nil {
			slog.Warn("closing agent", slog.String("error", err.Error()))
		}
	}()

	if configuration.ServerConfig.TLSConfig.CertificateFile != "" && configuration.ServerConfig.TLSConfig.PrivateKeyFile != "" {
		err = server.ListenAndServeTLS()
	} else {
		err = server.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
