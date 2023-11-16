package server

import (
	"github.com/fmiskovic/go-starter/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func initStaticRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(handlers.FlashMiddleware)

	app.Get("/", handlers.HandleHome)
	app.Get("/about", handlers.HandleAbout)
	app.Get("/flash", handlers.HandleFlash)

	app.Use(handlers.NotFoundMiddleware)
	app.Use(recover.New())
}
