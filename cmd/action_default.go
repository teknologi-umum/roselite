package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aldy505/sentry-integration/sloghandler"
	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/configor"
	slogmulti "github.com/samber/slog-multi"
	"github.com/teknologi-umum/roselite"
	"github.com/urfave/cli/v3"
)

func DefaultAction(ctx context.Context, c *cli.Command) error {
	var configuration Configuration
	err := configor.New(&configor.Config{}).Load(&configuration, c.String("config"))
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              configuration.ErrorReporting.SentryDSN,
		SampleRate:       configuration.ErrorReporting.SentrySampleRate,
		EnableTracing:    true,
		TracesSampleRate: configuration.ErrorReporting.SentryTracesSampleRate,
		Release:          version,
		EnableLogs:       true,
	})
	if err != nil {
		return fmt.Errorf("initializing Sentry: %w", err)
	}
	defer sentry.Flush(2 * time.Second)

	ctx = sentry.SetHubOnContext(ctx, sentry.CurrentHub())

	slog.SetDefault(slog.New(slogmulti.Fanout(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: configuration.ErrorReporting.LogLevel,
		}),
		sloghandler.NewSentrySlogHandler(configuration.ErrorReporting.LogLevel),
	)))

	monitors := make([]roselite.Monitor, len(configuration.Monitors))
	for i, monitor := range configuration.Monitors {
		slog.DebugContext(ctx, "Parsing monitor from configuration", slog.String("id", monitor.Id), slog.String("type", monitor.MonitorType), slog.String("target", monitor.MonitorTarget))
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
		EnableLogging:          true,
	})

	agent := roselite.NewAgent(roselite.AgentOptions{
		Monitors:               monitors,
		UpstreamKumaAddress:    configuration.UpstreamConfig.BaseUrl,
		UpstreamRequestHeaders: configuration.UpstreamConfig.RequestHeaders,
		UpstreamTLSConfig:      upstreamTLSConfig,
		RegionIdentifier:       configuration.Region,
		EnableLogging:          true,
	})

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal
		slog.InfoContext(ctx, "Received Ctrl+C. Shutting down server")
		err := server.Shutdown(ctx)
		if err != nil {
			slog.WarnContext(ctx, "closing server", slog.String("error", err.Error()))
		}

		slog.InfoContext(ctx, "Received Ctrl+C. Shutting down agent")
		err = agent.Close()
		if err != nil {
			slog.WarnContext(ctx, "closing agent", slog.String("error", err.Error()))
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			err := agent.Start()
			if err != nil {
				slog.WarnContext(ctx, "starting agent", slog.String("error", err.Error()))
			}

			slog.DebugContext(ctx, "agent process: sleeping for 5 seconds")
			time.Sleep(time.Second * 5)
		}
	}()

	if configuration.ServerConfig.TLSConfig.CertificateFile != "" && configuration.ServerConfig.TLSConfig.PrivateKeyFile != "" {
		slog.InfoContext(ctx, "Starting TLS Server, press Ctrl+C to exit.", slog.String("address", configuration.ServerConfig.ListenAddress))
		err = server.ListenAndServeTLS()
	} else {
		slog.InfoContext(ctx, "Starting Server, press Ctrl+C to exit.", slog.String("address", configuration.ServerConfig.ListenAddress))
		err = server.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}
