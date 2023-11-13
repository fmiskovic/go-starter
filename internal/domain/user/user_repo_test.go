package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/fmiskovic/go-starter/internal/containers"
	"github.com/fmiskovic/go-starter/internal/database"
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/fmiskovic/go-starter/util"
	"github.com/uptrace/bun"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestUserRepo_Create(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	// setup test-containers
	tearDown, ctx, bunDb := setUp(t)
	defer tearDown(t)

	repo := NewRepo(bunDb)

	// setup test cases
	type fields struct {
		repo UserRepo
	}
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:    "given valid user should not return error",
			fields:  fields{repo: repo},
			args:    args{u: createUser()},
			wantErr: nil,
		},
		{
			name:    "given nil user should return error",
			fields:  fields{repo: repo},
			args:    args{u: nil},
			wantErr: domain.NilEntityError,
		},
		{
			name:    "given user with id should return error",
			fields:  fields{repo: repo},
			args:    args{u: createUserWithId(1)},
			wantErr: errors.New("duplicate key value violates unique constraint"),
		},
		{
			name:    "given user with non-unique email should return error",
			fields:  fields{repo: repo},
			args:    args{u: createUser()},
			wantErr: errors.New("duplicate key value violates unique constraint \"users_email_key\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repo
			err := repo.Create(ctx, tt.args.u)

			if (err != nil || tt.wantErr != nil) && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("expected: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	// setup test-containers
	tearDown, ctx, bunDb := setUp(t)
	defer tearDown(t)

	repo := NewRepo(bunDb)

	// setup test cases
	type fields struct {
		repo UserRepo
	}
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		given   func() error
		verify  func(t *testing.T)
		wantErr error
	}{
		{
			name:   "given valid user input should not return error",
			fields: fields{repo: repo},
			args:   args{u: createUserWithId(1, WithEmail("test@test.com"))},
			given: func() error {
				return repo.Create(ctx, createUser())
			},
			verify: func(t *testing.T) {
				u, err := repo.GetById(ctx, 1)
				if err != nil {
					t.Errorf("failed to get user by id 1, error: %v", err)
				}
				if "test@test.com" != u.Email {
					t.Errorf("expected: %v, got: %v", "test@test.com", u.Email)
				}
			},
			wantErr: nil,
		},
		{
			name:   "given nil user input should return error",
			fields: fields{repo: repo},
			args:   args{u: nil},
			given: func() error {
				return nil
			},
			verify:  func(t *testing.T) {},
			wantErr: domain.NilEntityError,
		},
		{
			name:   "given user with non existing id should return error",
			fields: fields{repo: repo},
			args:   args{u: createUserWithId(111)},
			given: func() error {
				return nil
			},
			verify: func(t *testing.T) {
				_, err := repo.GetById(ctx, 111)

				if err == nil || !strings.Contains(err.Error(), "no rows in result set") {
					t.Errorf("expected: %v, got: %v", "no rows in result set", err)
				}
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.given(); err != nil {
				t.Errorf("failed to execute given func: %v", err)
			}

			err := repo.Update(ctx, tt.args.u)
			if (err != nil || tt.wantErr != nil) && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Errorf("expected: %v, got: %v", tt.wantErr, err)
				return
			}

			tt.verify(t)
		})
	}
}

func setUp(t *testing.T) (func(t *testing.T), context.Context, *bun.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	// start postgres container
	postgres, err := containers.StartPostgresContainer(ctx)
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
				database.ConnString,
				dbUser,
				dbPassword,
				fmt.Sprintf("%s:%d", host, port.Int()),
				dbName,
			)

			bunDb = database.Connect(dbUri)

			if err := containers.MigrateDB(ctx, bunDb); err != nil {
				t.Fatalf("db migration failed: %v", err)
			}
			break
		}
		slog.Info("waiting for postgres container...")
	}

	return func(t *testing.T) {
		if err := containers.TerminateContainer(ctx, postgres); err != nil {
			slog.Warn("failed to terminate container", err)
		}
		cancel()
	}, ctx, bunDb
}

func createUserWithId(id int64, opts ...Option) *User {
	u := createUser()
	u.ID = id
	for _, opt := range opts {
		opt(u)
	}
	return u
}
