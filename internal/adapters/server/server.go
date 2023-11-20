// Package server represents secondary adapter.
package server

import (
	"errors"
	"github.com/fmiskovic/go-starter/internal/adapters/repos"
	"github.com/fmiskovic/go-starter/internal/adapters/server/config"
	"github.com/fmiskovic/go-starter/internal/adapters/web/handlers"
	"github.com/fmiskovic/go-starter/internal/adapters/web/routes"
	"github.com/fmiskovic/go-starter/internal/ports"
	"github.com/gofiber/template/django/v3"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

// Server holds configuration, database connection and fiber app.
type Server struct {
	Config config.ServerConfig
	Db     *bun.DB
	App    *fiber.App
}

// New instantiate new Server with specified config.
func New(config config.ServerConfig) *Server {
	return &Server{Config: config}
}

// InitDb connects Server to the DB.
func (s *Server) InitDb() error {
	s.Db = ports.Connect(s.Config.DbConnString, s.Config.MaxOpenConn, s.Config.MaxIdleConn)
	return nil
}

// InitApp initialises fiber app.
func (s *Server) InitApp() error {
	if s.Db == nil {
		return errors.New("DB must be initialized first")
	}
	app := fiber.New(fiber.Config{
		ErrorHandler:          handlers.ErrorHandler,
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 initViews(),
	})

	// init swagger
	routes.InitSwaggerRoutes(app)

	// init user api handlers
	routes.InitUserRoutes(repos.NewUserRepo(s.Db), handlers.NewValidator(), app)
	// init static handlers
	routes.InitStaticRoutes(app)

	s.App = app
	return nil
}

// Ready returns true if everything is properly configured.
func (s *Server) Ready() bool {
	return s.Db != nil && s.App != nil
}

// Start the server.
func (s *Server) Start() error {
	if !s.Ready() {
		return errors.New("server is not ready")
	}

	slog.Info("the app is up and running...", "address", s.Config.ListenAddr)
	return s.App.Listen(s.Config.ListenAddr)
}

func initViews() *django.Engine {
	engine := django.New("./views", ".html")
	engine.Reload(true)
	engine.AddFunc("css", func(name string) (res template.HTML) {
		err := filepath.Walk("public/assets", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == name {
				res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
			}
			return nil
		})
		if err != nil {
			log.Fatalf("failed to create django engine, unable to walk puiblic/assets folder. Error: %v", err)
		}
		return
	})
	return engine
}
