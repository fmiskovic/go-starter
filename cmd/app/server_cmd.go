package main

import (
	"github.com/fmiskovic/go-starter/internal/config"
	"github.com/fmiskovic/go-starter/internal/server"
	"github.com/urfave/cli/v2"
	"log/slog"
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

			err := srv.Start()
			if err == nil {
				slog.Info("the app is up and running...", "address", srv.Config.ListenAddr)
			}
			return err
		},
	}
}
