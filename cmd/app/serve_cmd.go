package main

import (
	"github.com/urfave/cli/v2"
)

// newServeCmd configures start server cli command.
func newServeCmd() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "start the server",
		Action: func(ctx *cli.Context) error {
			return newServer(defaultConfig).start()
		},
	}
}
