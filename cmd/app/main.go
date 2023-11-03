package main

import (
	"github.com/fmiskovic/go-starter/migrations"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "app",

		Commands: []*cli.Command{
			newServeCmd(),
			newMigrationCmd(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
