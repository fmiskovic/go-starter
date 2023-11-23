package api

import (
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserRouter struct {
	repo ports.UserRepo[uuid.UUID]
	app  *fiber.App
}

func NewUserRouter(repo ports.UserRepo[uuid.UUID], app *fiber.App) UserRouter {
	return UserRouter{repo: repo, app: app}
}

// InitRouters initializes user management api.
func (r UserRouter) InitRouters() {
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	handler := NewUserHandler(r.repo)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}
