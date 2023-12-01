package main

import (
	"github.com/fmiskovic/go-starter/internal/adapters/handlers/user"
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/adapters/views"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type UserRouter struct {
	service    services.UserService
	app        *fiber.App
	authConfig configs.AuthConfig
}

// NewRouter instantiates new user.Router
func NewUserRouter(db *bun.DB, app *fiber.App, authConfig configs.AuthConfig) UserRouter {
	repo := repos.NewUserRepo(db)
	svc := services.NewUserService(repo, authConfig)
	return UserRouter{service: svc, app: app, authConfig: authConfig}
}

// InitRouters initializes user management api.
func (r UserRouter) InitRouters() {
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	handler := user.NewHandler(r.service)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}

// InitStaticRoutes initializes static view handlers to serve the UI.
func initStaticRouters(app *fiber.App) {
	app.Static("/public", "./public")

	app.Use(views.FlashMiddleware)

	app.Get("/", views.HandleHome)
	app.Get("/about", views.HandleAbout)
	app.Get("/flash", views.HandleFlash)

	app.Use(views.NotFoundMiddleware)
}

// InitSwaggerRoutes initializes Swagger UI.
func initSwaggerRouters(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.json",
		Path:     "docs",
		Title:    "API Documentation",
	}))
}

// Protect is auth middleware used to protect endpoints.
func Protect(secret string) func(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}
