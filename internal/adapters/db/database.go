package db

import (
	"database/sql"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Database holds properties needed for making db connection.
type Database struct {
	Uri         string
	MaxOpenConn int
	MaxIdleConn int
}

func NewDatabase(uri string, maxOpen int, maxIdle int) Database {
	return Database{Uri: uri, MaxOpenConn: maxOpen, MaxIdleConn: maxIdle}
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
	// if utils.IsDev() {
	// 	bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	// }

	return bunDb, nil
}
