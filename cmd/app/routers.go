package main

import (
	"github.com/fmiskovic/go-starter/internal/adapters/handlers"
	"github.com/fmiskovic/go-starter/internal/adapters/handlers/auth"
	"github.com/fmiskovic/go-starter/internal/adapters/handlers/user"
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/core/configs"
	"github.com/fmiskovic/go-starter/internal/core/services"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Router struct {
	service        services.UserService
	app            *fiber.App
	authConfig     configs.AuthConfig
	authMiddleware auth.Middleware
}

// NewRouter instantiates new user.Router
func newRouter(db *bun.DB, app *fiber.App, authConfig configs.AuthConfig) Router {
	repo := repos.NewUserRepo(db)
	svc := services.NewUserService(repo, authConfig)
	authMiddleware := auth.NewMiddleware(authConfig)
	return Router{service: svc, app: app, authConfig: authConfig, authMiddleware: authMiddleware}
}

// initUserRouters initializes user management api.
func (r Router) initUserRouters() {
	api := r.app.Group("/api")
	v1 := api.Group("/v1")
	userGroup := v1.Group("/user", r.authMiddleware.AdminAuthenticated())

	handler := user.NewHandler(r.service)

	userGroup.Get("/:id", handler.HandleGetById())
	userGroup.Get("/", handler.HandleGetPage())
	userGroup.Delete("/:id", handler.HandleDeleteById())
	userGroup.Post("/", handler.HandleCreate())
	userGroup.Put("/", handler.HandleUpdate())
	userGroup.Post("/roles", handler.HandleUserRoles())
	userGroup.Post("/:id/enabledisable", handler.HandleEnableDisable())
}

func (r Router) initAuthRouters() {
	a := r.app.Group("/auth")

	handler := auth.NewHandler(r.service)
	a.Post("/login", handler.HandleSignIn())
	a.Get("/logout", handler.HandleSignOut())
	a.Post("/register", handler.HandleSignUp())
	a.Post("/email", handler.HandleConfirmEmail())
	a.Post("/password", handler.HandleChangePassword())
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
