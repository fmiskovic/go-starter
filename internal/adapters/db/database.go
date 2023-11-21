package db

import (
	"database/sql"
	"github.com/fmiskovic/go-starter/internal/utils"
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

// OpenDb implements ports.Db interface and opens new connection.
// Returns *bun.DB, or error if connection failed.
func (db Database) OpenDb() (*bun.DB, error) {
	slog.Info("initializing db with conn string", "conn", db.Uri)

	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(db.Uri)))
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}

	sqlDb.SetMaxOpenConns(db.MaxOpenConn)
	sqlDb.SetMaxIdleConns(db.MaxIdleConn)
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if utils.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return bunDb, nil
}
