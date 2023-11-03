package user

import "github.com/gofiber/fiber/v2"

func InitRoutes(repo UserRepo, app *fiber.App) {
	app.Get("/user/:id", HandleGetById(repo))
	app.Delete("/user/:id", HandleDeleteById(repo))
	app.Post("/user", HandleCreate(repo))
	app.Put("/user", HandleUpdate(repo))
}
