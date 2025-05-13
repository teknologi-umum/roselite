package main

import (
	"context"
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/urfave/cli/v3"
)

func AgentAction(ctx context.Context, c *cli.Command) error {
	var configuration Configuration
	err := configor.New(&configor.Config{}).Load(&configuration, c.String("config"))
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	return nil
}
