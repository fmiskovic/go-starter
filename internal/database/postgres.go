package database

import (
	"database/sql"
	"fmt"
	"github.com/fmiskovic/go-starter/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log/slog"
	"runtime"
	"sync"
)

var DbBun *bun.DB

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
		maxOpenConns := runtime.NumCPU() + 1
		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxOpenConns)
		DbBun = bun.NewDB(sqldb, pgdialect.New())
		if util.IsDev() {
			DbBun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
		}
	})
}

func Connect(uri string) *bun.DB {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if util.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return bunDb
}
