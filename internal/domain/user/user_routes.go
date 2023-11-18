package user

import (
	"github.com/fmiskovic/go-starter/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(repo UserRepo, validator validator.Validator, app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/user/:id", HandleGetById(repo))
	v1.Get("/user", HandleGetPage(repo))
	v1.Delete("/user/:id", HandleDeleteById(repo))
	v1.Post("/user", HandleCreate(repo, validator))
	v1.Put("/user", HandleUpdate(repo))
}
