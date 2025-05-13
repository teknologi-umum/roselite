package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v3"
)

var version string
var copyright = `Copyright (C) 2025  Teknologi Umum <opensource@teknologiumum.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

func main() {
	cmd := &cli.Command{
		Name:           "roselite",
		Aliases:        []string{},
		Usage:          "Active relay for Uptime Kuma's push monitor type",
		Version:        version,
		Description:    "Active relay for Uptime Kuma's push monitor type. It is also compatible with Semyi (made by Teknologi Umum), in which more data related to the monitor is sent to Semyi instance.",
		DefaultCommand: "agent",
		Commands: []*cli.Command{
			{
				Name:    "agent",
				Version: version,
				Usage:   "Start Roselite in agent mode, it will not expose any HTTP port",
				Action: func(ctx context.Context, c *cli.Command) error {
					return nil
				},
			},
			{
				Name:    "server",
				Version: version,
				Usage:   "Start Roselite in server mode, it will expose HTTP port",
				Action: func(ctx context.Context, c *cli.Command) error {
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Path to Roselite configuration file",
				HideDefault: false,
				Sources:     cli.ValueSourceChain{Chain: []cli.ValueSource{cli.EnvVar("CONFIGURATION_FILE_PATH")}},
				Required:    true,
				OnlyOnce:    true,
			},
		},
		Authors:   []any{},
		Copyright: copyright,
		Suggest:   true,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
