package main

import (
    "context"
    "fmt"
    "log/slog"
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

func AgentAction(ctx context.Context, c *cli.Command) error {
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
        slog.InfoContext(ctx, "Received Ctrl+C, closing agent.")
        err := agent.Close()
        if err != nil {
            slog.WarnContext(ctx, "closing agent", slog.String("error", err.Error()))
        }
    }()

    slog.InfoContext(ctx, "Starting agent, press Ctrl+C to exit.")
    if err := agent.Start(); err != nil {
        return fmt.Errorf("starting agent: %w", err)
    }
    return nil
}
