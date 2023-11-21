package main

import (
	"github.com/fmiskovic/go-starter/internal/adapters/views"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// InitStaticRoutes initializes static view handlers to serve the UI.
func initStaticRoutes(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(views.FlashMiddleware)

	app.Get("/", views.HandleHome)
	app.Get("/about", views.HandleAbout)
	app.Get("/flash", views.HandleFlash)

	app.Use(views.NotFoundMiddleware)
}

// InitSwaggerRoutes initializes Swagger UI.
func initSwaggerRoutes(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "API Documentation",
	}))
}
