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

// Database holds properties needed for making db connection.
type Database struct {
	Uri         string
	MaxOpenConn int
	MaxIdleConn int
}

// Connect opens new DB connection with specified configuration.
func (db Database) Connect() *bun.DB {
	slog.Info("initializing db with conn string", "conn", db.Uri)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(db.Uri)))

	sqlDb.SetMaxOpenConns(db.MaxOpenConn)
	sqlDb.SetMaxIdleConns(db.MaxIdleConn)
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if helper.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return bunDb
}
