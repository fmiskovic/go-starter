package api

import (
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/gofiber/fiber/v2"
)

// InitUserRoutes initializes user management api.
func InitUserRoutes(repo ports.UserRepo[uint64], app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	handler := NewUserHandler(repo)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}
