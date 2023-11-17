package test

import (
	"context"
	"fmt"
	"github.com/fmiskovic/go-starter/internal/database"
	"github.com/fmiskovic/go-starter/migrations"
	"github.com/fmiskovic/go-starter/util"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"log/slog"
	"runtime"
	"testing"
	"time"
)

func SetUp(t *testing.T) (func(t *testing.T), context.Context, *bun.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	// start postgres container
	postgres, err := startPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	var bunDb *bun.DB
	for {
		if postgres.IsRunning() {
			slog.Info("postgres container is ready")
			host, err := postgres.Host(ctx)
			if err != nil {
				panic(err)
			}

			port, err := postgres.MappedPort(ctx, "5432")
			if err != nil {
				panic(err)
			}

			dbName := util.GetEnvOrDefault("DB_NAME", "test-db")
			dbUser := util.GetEnvOrDefault("DB_USER", "test")
			dbPassword := util.GetEnvOrDefault("DB_PASSWORD", "test")

			dbUri := fmt.Sprintf(
				"postgresql://%s:%s@%s/%s?sslmode=disable",
				dbUser,
				dbPassword,
				fmt.Sprintf("%s:%d", host, port.Int()),
				dbName,
			)

			openDbConn := runtime.NumCPU() + 1
			bunDb = database.Connect(dbUri, openDbConn, openDbConn)

			if err := migrateDB(ctx, bunDb); err != nil {
				t.Fatalf("db migration failed: %v", err)
			}
			break
		}
		slog.Info("waiting for postgres container...")
	}

	return func(t *testing.T) {
		if err := terminateContainer(ctx, postgres); err != nil {
			slog.Warn("failed to terminate container", err)
		}
		cancel()
	}, ctx, bunDb
}

func startPostgresContainer(ctx context.Context) (testcontainers.Container, error) {
	dbName := util.GetEnvOrDefault("DB_NAME", "test-db")
	dbUser := util.GetEnvOrDefault("DB_USER", "test")
	dbPassword := util.GetEnvOrDefault("DB_PASSWORD", "test")

	// Define a Postgres container configuration.
	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	// create and start the postgres container.
	return testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
}

func migrateDB(ctx context.Context, db *bun.DB) error {
	slog.Info("db migration is starting...")
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("init failed: %v", err)
	}

	if err := migrator.Lock(ctx); err != nil {
		slog.Warn("lock failed but it's ok, error message:", err)
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	if group.IsZero() {
		fmt.Printf("there are no new migrations to run (database is up to date)\n")
	}

	fmt.Printf("migrated to %s\n", group)

	return nil
}

func terminateContainer(ctx context.Context, container testcontainers.Container) error {
	if container == nil {
		slog.Info("container is nil, skipping terminate func")
		return nil
	}
	return container.Terminate(ctx)
}
