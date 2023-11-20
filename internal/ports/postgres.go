// Package database is a port to connect to a db.
package ports

import (
	"database/sql"
	"github.com/fmiskovic/go-starter/internal/helper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log/slog"
)

// Connect opens new DB connection with specified uri, maxOpenConn, maxIdleConn, and returns pointer to a bun.DB.
func Connect(uri string, maxOpenConn int, maxIdleConn int) *bun.DB {
	slog.Info("initializing db with conn string", "conn", uri)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	sqlDb.SetMaxOpenConns(maxOpenConn)
	sqlDb.SetMaxIdleConns(maxIdleConn)
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if helper.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return bunDb
}
