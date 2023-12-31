package main

import (
	"errors"
	"html/template"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/fmiskovic/go-starter/internal/adapters/db"
	"github.com/fmiskovic/go-starter/internal/utils"

	"github.com/fmiskovic/go-starter/internal/core/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django/v3"
	"github.com/uptrace/bun"
)

// Server holds configuration, database connection and fiber app.
type Server struct {
	Config ServerConfig
	Db     *bun.DB
	App    *fiber.App
}

// newServer instantiate new Server with specified config.
func newServer(config ServerConfig) Server {
	bunDb, err := initDb(config)
	if err != nil {
		log.Fatal(err)
	}
	app := initApp(bunDb, config.AuthConfig)
	return Server{
		Config: config,
		Db:     bunDb,
		App:    app,
	}
}

// Ready returns true if everything is properly configured.
func (s Server) ready() bool {
	return s.Db != nil && s.App != nil
}

// Start the server.
func (s Server) start() error {
	if !s.ready() {
		return errors.New("server is not ready")
	}

	slog.Info("the app is up and running...", "address", s.Config.ListenAddr)
	return s.App.Listen(s.Config.ListenAddr)
}

// ----- INITS ----- //

func initDb(config ServerConfig) (*bun.DB, error) {
	return db.Database{
		Uri:         config.DbConnString,
		MaxOpenConn: config.MaxOpenConn,
		MaxIdleConn: config.MaxOpenConn,
	}.OpenDb()
}

func initApp(db *bun.DB, authConfig configs.AuthConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		PassLocalsToViews:     true,
		Views:                 initViews(),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: utils.GetEnvOrDefault("ALLOW_ORIGINS", "*"),
		MaxAge:       -1, //negative number disables caching completely
	}))

	if utils.IsDev() {
		app.Use(pprof.New())
	}

	router := newRouter(db, app, authConfig)

	// init swagger
	router.initSwaggerRouters()
	// init auth routers
	router.initAuthRouters()
	// init user api handlers
	router.initUserRouters()
	// init static handlers
	router.initStaticRouters()

	app.Use(recover.New())
	return app
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
