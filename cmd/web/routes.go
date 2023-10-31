package main

import (
	"github.com/fmiskovic/go-starter/handlers"
	"github.com/gofiber/fiber/v2"
)

func initRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(handlers.FlashMiddleware)

	app.Get("/", handlers.HandleHome)
	app.Get("/about", handlers.HandleAbout)
	app.Get("/flash", handlers.HandleFlash)

	app.Use(handlers.NotFoundMiddleware)
}
