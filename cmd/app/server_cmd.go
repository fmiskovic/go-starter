package main

import (
	"github.com/fmiskovic/go-starter/internal/config"
	"github.com/fmiskovic/go-starter/internal/server"
	"github.com/urfave/cli/v2"
)

func newServeCmd() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "start the server",
		Action: func(ctx *cli.Context) error {
			srv := server.New(config.DefaultConfig)
			if err := srv.InitDb(); err != nil {
				return err
			}
			if err := srv.InitApp(); err != nil {
				return err
			}
			return srv.Start()
		},
	}
}
