package routes

import (
	"github.com/fmiskovic/go-starter/internal/adapters/web/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// InitStaticRoutes initializes static view handlers to serve the UI.
func InitStaticRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(handlers.FlashMiddleware)

	app.Get("/", handlers.HandleHome)
	app.Get("/about", handlers.HandleAbout)
	app.Get("/flash", handlers.HandleFlash)

	app.Use(handlers.NotFoundMiddleware)
	app.Use(recover.New())
}
