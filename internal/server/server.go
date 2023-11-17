package server

import (
	"errors"
	"github.com/fmiskovic/go-starter/internal/config"
	"github.com/fmiskovic/go-starter/internal/database"
	"github.com/fmiskovic/go-starter/internal/domain/user"
	"github.com/fmiskovic/go-starter/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
)

type Server struct {
	Config config.AppConfig
	Db     *bun.DB
	App    *fiber.App
}

func New(config config.AppConfig) *Server {
	return &Server{Config: config}
}

func (s *Server) InitDb() error {
	s.Db = database.Connect(s.Config.DbConnString, s.Config.MaxOpenConn, s.Config.MaxIdleConn)
	return nil
}

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

	// init user api handlers
	user.InitRoutes(user.NewRepo(s.Db), app)
	// init static handlers
	initStaticRoutes(app)

	s.App = app
	return nil
}

func (s *Server) Ready() bool {
	return s.Db != nil && s.App != nil
}

func (s *Server) Start() error {
	if !s.Ready() {
		return errors.New("server is not ready")
	}

	return s.App.Listen(s.Config.ListenAddr)
}
