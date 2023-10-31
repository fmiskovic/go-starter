package user

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fmiskovic/go-starter/database"
	"github.com/fmiskovic/go-starter/migrations"
	"github.com/fmiskovic/go-starter/util"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun/migrate"
)

const Timeout = time.Second * 30

func TestUserRepo(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	// ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	// defer cancel()

	ctx := context.Background()

	container := runPostgres(ctx)

	// Clean up the container
	defer terminatePostgres(ctx, container)

	setEnvs(ctx, t, container)
	runMigration(ctx, container)

	repo := NewRepo()

	u := createUser()

	err := repo.Save(ctx, u)
	if err != nil {
		t.Errorf("test failed %v", err)
	}

}

// runPostgres starts a PostgreSQL container and returns a handle to it.
func runPostgres(ctx context.Context) testcontainers.Container {
	dbName := util.GetEnvOrDefault("DB_NAME", "go-database")
	dbUser := util.GetEnvOrDefault("DB_USER", "dbadmin")
	dbPassword := util.GetEnvOrDefault("DB_PASSWORD", "dbadmin")

	// Define a PostgreSQL container configuration.
	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	// Create and start the PostgreSQL container.
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		panic(err)
	}

	return container
}

func terminatePostgres(ctx context.Context, container testcontainers.Container) {
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}
}

func setEnvs(ctx context.Context, t *testing.T, container testcontainers.Container) {
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	os.Setenv("DB_HOST", fmt.Sprintf("%s:%d", host, port.Int()))
	// t.Setenv("DB_HOST", fmt.Sprintf("%s:%d", host, port.Int()))
}

func runMigration(ctx context.Context, container testcontainers.Container) {
	// p, err := container.MappedPort(context.Background(), "5432")
	// if err != nil {
	// 	panic(err)
	// }

	// h, err := container.Host(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// user := util.GetEnvOrDefault("DB_USER", "dbadmin")
	// password := util.GetEnvOrDefault("DB_PASSWORD", "dbadmin")
	// host := util.GetEnvOrDefault("DB_HOST", fmt.Sprintf("%s:%d", h, p.Int()))
	// name := util.GetEnvOrDefault("DB_NAME", "go-database")
	// uri := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
	// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))

	// database := bun.NewDB(sqldb, pgdialect.New())

	migrator := migrate.NewMigrator(database.Bun, migrations.Migrations)

	defer migrator.Unlock(ctx) //nolint:errcheck

	if err := migrator.Init(ctx); err != nil {
		panic(err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		panic(err)
	}
	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
	}
	fmt.Printf("migrated to %s\n", group)
}
