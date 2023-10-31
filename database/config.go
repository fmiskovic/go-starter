package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fmiskovic/go-starter/util"
	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var Bun *bun.DB

const Timeout = time.Second * 10

type Config struct {
	User     string
	Password string
	Host     string
	Name     string
}

func init() {
	if err := util.LoadEnvVars(); err != nil {
		slog.Warn("unable to locate .env file, default environment values will be used")
	}

	var (
		user     = util.GetEnvOrDefault("DB_USER", "dbadmin")
		password = util.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
		host     = util.GetEnvOrDefault("DB_HOST", "localhost:5432")
		name     = util.GetEnvOrDefault("DB_NAME", "go-database")
		uri      = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
		sqldb    = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
		once     = sync.Once{}
	)
	fmt.Println(uri)
	once.Do(func() {
		Bun = bun.NewDB(sqldb, pgdialect.New())
		if util.IsDev() {
			Bun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}
	})
}
