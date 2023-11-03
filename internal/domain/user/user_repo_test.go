package user

import (
	"context"
	"fmt"
	database2 "github.com/fmiskovic/go-starter/internal/database"
	"github.com/fmiskovic/go-starter/migrations"
	"github.com/uptrace/bun"
	"testing"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	container := runPostgres(ctx)
	defer terminatePostgres(ctx, container)

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	dbUri := fmt.Sprintf(
		database2.ConnString,
		"dbadmin",
		"dbadmin",
		fmt.Sprintf("%s:%d", host, port.Int()),
		"go-db",
	)

	bunDb := database2.Connect(dbUri)

	runMigration(ctx, bunDb)

	repo := NewRepo(bunDb)

	u := createUser()

	err = repo.Save(ctx, u)
	if err != nil {
		t.Errorf("test failed %v", err)
	}

}

// runPostgres starts a Postgres container and returns a handle to it.
func runPostgres(ctx context.Context) testcontainers.Container {
	dbName := util.GetEnvOrDefault("DB_NAME", "go-db")
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

func runMigration(ctx context.Context, db *bun.DB) {
	migrator := migrate.NewMigrator(db, migrations.Migrations)

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
