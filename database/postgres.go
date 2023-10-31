package database

import (
	"database/sql"
	"fmt"
	"github.com/fmiskovic/go-starter/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresDb struct {
	BunDB *bun.DB
}

func Connect(uri string) *bun.DB {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	bunDb := bun.NewDB(sqlDb, pgdialect.New())
	if util.IsDev() {
		bunDb.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	return bunDb
}

func ConnectDefault() *bun.DB {
	uri := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		DefaultConfig.User,
		DefaultConfig.Password,
		DefaultConfig.Host,
		DefaultConfig.Name,
	)

	return Connect(uri)
}
