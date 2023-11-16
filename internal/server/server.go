package server

import (
	"database/sql"
	"errors"
	"github.com/fmiskovic/go-starter/internal/config"
	"github.com/fmiskovic/go-starter/internal/database"
	"github.com/fmiskovic/go-starter/internal/domain/user"
	"github.com/fmiskovic/go-starter/internal/handlers"
	"github.com/fmiskovic/go-starter/util"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log/slog"
)

type Server struct {
	Config config.ServerConfig
	Db     *bun.DB
	App    *fiber.App
}

func NewServer(config config.ServerConfig) *Server {
	return &Server{Config: config}
}

func InitDb(s *Server) error {
	dbConnString := s.Config.DbConnString
	slog.Info("initializing db with conn string", "conn", dbConnString)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbConnString)))

	sqlDb.SetMaxOpenConns(s.Config.MaxOpenConn)
	sqlDb.SetMaxIdleConns(s.Config.MaxIdleConn)
	db := bun.NewDB(sqlDb, pgdialect.New())
	if util.IsDev() {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	s.Db = db

	return nil
}

func InitApp(s *Server) error {
	app := fiber.New(fiber.Config{
		ErrorHandler:          handlers.ErrorHandler,
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 initViews(),
	})

	// init user api handlers
	user.InitRoutes(user.NewRepo(database.DbBun), app)
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
