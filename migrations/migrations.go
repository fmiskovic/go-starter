package migrations

import (
	"embed"
	"fmt"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	fmt.Println("initializing migrrations...")

	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
