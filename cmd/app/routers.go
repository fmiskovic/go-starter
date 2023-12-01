package main

import (
	"github.com/fmiskovic/go-starter/internal/adapters/handlers"
	"github.com/fmiskovic/go-starter/internal/adapters/handlers/user"
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Router struct {
	service    services.UserService
	app        *fiber.App
	authConfig configs.AuthConfig
}

// NewRouter instantiates new user.Router
func newRouter(db *bun.DB, app *fiber.App, authConfig configs.AuthConfig) Router {
	repo := repos.NewUserRepo(db)
	svc := services.NewUserService(repo, authConfig)
	return Router{service: svc, app: app, authConfig: authConfig}
}

// initUserRouters initializes user management api.
func (r Router) initUserRouters() {
	api := r.app.Group("/api")
	v1 := api.Group("/v1")

	handler := user.NewHandler(r.service)

	v1.Get("/user/:id", handler.HandleGetById())
	v1.Get("/user", handler.HandleGetPage())
	v1.Delete("/user/:id", handler.HandleDeleteById())
	v1.Post("/user", handler.HandleCreate())
	v1.Put("/user", handler.HandleUpdate())
}

// initStaticRoutes initializes static view handlers to serve the UI.
func (r Router) initStaticRouters() {
	r.app.Static("/public", "./public")

	r.app.Use(handlers.FlashMiddleware)

	r.app.Get("/", handlers.HandleHome)
	r.app.Get("/about", handlers.HandleAbout)
	r.app.Get("/flash", handlers.HandleFlash)

	r.app.Use(handlers.NotFoundMiddleware)
}

// initSwaggerRoutes initializes Swagger UI.
func (r Router) initSwaggerRouters() {
	r.app.Use(swagger.New(swagger.Config{
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
