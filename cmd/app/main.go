package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fmiskovic/bs/handlers"
	"github.com/fmiskovic/bs/util"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:          handlers.ErrorHandler,
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 createEngine(),
	})

	initRoutes(app)
	listenAddr := listenAddrOrDefault()
	fmt.Printf("app running in %s and listening on: %s\n", util.AppEnv(), listenAddr)
	log.Fatal(app.Listen(listenAddr))
}

func initRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(handlers.FlashMiddleware)

	app.Get("/", handlers.HandleHome)
	app.Get("/bored", handlers.HandleBored)
	app.Get("/flash", handlers.HandleFlash)

	app.Use(handlers.NotFoundMiddleware)
}

func listenAddrOrDefault() string {
	addr := os.Getenv("HTTP_LISTEN_ADDR")
	if strings.TrimSpace(addr) == "" {
		return ":3000"
	}
	return addr
}
