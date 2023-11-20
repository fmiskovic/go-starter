package routes

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// InitSwaggerRoutes initializes Swagger UI.
func InitSwaggerRoutes(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "API Documentation",
	}))
}
