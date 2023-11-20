package routes

import (
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/adapters/web/handlers"
	"github.com/gofiber/fiber/v2"
)

// InitUserRoutes initializes user management endpoints.
func InitUserRoutes(repo repos.UserRepo, validator handlers.Validator, app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/user/:id", handlers.HandleGetById(repo))
	v1.Get("/user", handlers.HandleGetPage(repo))
	v1.Delete("/user/:id", handlers.HandleDeleteById(repo))
	v1.Post("/user", handlers.HandleCreate(repo, validator))
	v1.Put("/user", handlers.HandleUpdate(repo, validator))
}
