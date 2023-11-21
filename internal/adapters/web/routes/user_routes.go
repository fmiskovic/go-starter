package routes

import (
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/adapters/web/handlers"
	"github.com/gofiber/fiber/v2"
)

// InitUserRoutes initializes user management endpoints.
func InitUserRoutes(repo repos.UserRepo, app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	handler := handlers.NewUserHandler(repo)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}
