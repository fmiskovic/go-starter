package routes

import (
	handlers2 "github.com/fmiskovic/go-starter/internal/interfaces/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// InitStaticRoutes initializes static view handlers to serve the UI.
func InitStaticRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(handlers2.FlashMiddleware)

	app.Get("/", handlers2.HandleHome)
	app.Get("/about", handlers2.HandleAbout)
	app.Get("/flash", handlers2.HandleFlash)

	app.Use(handlers2.NotFoundMiddleware)
	app.Use(recover.New())
}
