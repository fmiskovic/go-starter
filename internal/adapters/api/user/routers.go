package user

import (
	"github.com/fmiskovic/go-starter/internal/core/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Router struct {
	repo ports.UserRepo[uuid.UUID]
	app  *fiber.App
}

// NewRouter instantiates new user.Router
func NewRouter(repo ports.UserRepo[uuid.UUID], app *fiber.App) Router {
	return Router{repo: repo, app: app}
}

// InitRouters initializes user management api.
func (r Router) InitRouters() {
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	handler := NewHandler(r.repo)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}
