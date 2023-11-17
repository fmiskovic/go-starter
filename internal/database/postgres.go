package database

import (
	"database/sql"
	"github.com/fmiskovic/go-starter/pkg/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log/slog"
)

func Connect(uri string, maxOpenConn int, maxIdleConn int) *bun.DB {
	slog.Info("initializing db with conn string", "conn", uri)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	sqlDb.SetMaxOpenConns(maxOpenConn)
	sqlDb.SetMaxIdleConns(maxIdleConn)
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if util.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return bunDb
}
