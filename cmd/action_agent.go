package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jinzhu/configor"
	"github.com/teknologi-umum/roselite"
	"github.com/urfave/cli/v3"
)

func AgentAction(ctx context.Context, c *cli.Command) error {
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

	agent := roselite.NewAgent(roselite.AgentOptions{
		Monitors:               monitors,
		UpstreamKumaAddress:    configuration.UpstreamConfig.BaseUrl,
		UpstreamRequestHeaders: configuration.UpstreamConfig.RequestHeaders,
		UpstreamTLSConfig:      upstreamTLSConfig,
		RegionIdentifier:       configuration.Region,
	})

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)

	go func() {
		<-exitSignal
		err := agent.Close()
		if err != nil {
			slog.Warn("closing agent", slog.String("error", err.Error()))
		}
	}()

	if err := agent.Start(); err != nil {
		return fmt.Errorf("starting agent: %w", err)
	}
	return nil
}
