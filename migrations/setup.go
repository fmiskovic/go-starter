// Package migrations: You should put each migration into a separate sql file.
// A migration file names consists of an unique migration name (20210505110026) and a comment (add_foo_column),
// for example, 20210505110026_add_foo_column.up.sql.
// For more info, read https://bun.uptrace.dev/guide/migrations.html#sql-based-migrations
package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

// Migrations object.
var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
