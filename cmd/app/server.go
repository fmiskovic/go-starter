package main

import (
	"fmt"
	"github.com/fmiskovic/go-starter/internal/handlers"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/fmiskovic/go-starter/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"github.com/urfave/cli/v2"
)

func newServeCmd() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "start API server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "addr",
				Value: ":8080",
				Usage: "serve address",
			},
		},
		Action: func(ctx *cli.Context) error {
			app := initApp()
			listenAddr := listenAddrOrDefault(ctx)
			fmt.Printf("app is running in %s environment and listening on: %s\n", util.AppEnv(), listenAddr)
			return app.Listen(listenAddr)
		},
	}
}

func listenAddrOrDefault(ctx *cli.Context) string {
	addr := ctx.String("addr")
	if util.IsBlank(addr) {
		addr = util.GetEnvOrDefault("HTTP_LISTEN_ADDR", ":8080")
	}
	return addr
}

func initApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          handlers.ErrorHandler,
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 initViews(),
	})
	initRoutes(app)
	return app
}

func initViews() *django.Engine {
	engine := django.New("./views", ".html")
	engine.Reload(true)
	engine.AddFunc("css", func(name string) (res template.HTML) {
		err := filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
			}
			return nil
		})
		if err != nil {
			log.Fatalf("failed to create django engine, unable to walk puiblic/assets folder. Error: %v", err)
		}
		return
	})
	return engine
}
