package user

import (
	"errors"
	"github.com/fmiskovic/go-starter/internal/containers"
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/matryer/is"
	"strings"
	"testing"
	"time"
)

func TestUserRepo_GetById(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.NewRelaxed(t)

	// setup test-containers
	tearDown, ctx, bunDb := containers.SetUp(t)
	defer tearDown(t)

	repo := NewRepo(bunDb)

	// setup test cases
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		given   func(t *testing.T) error
		wantErr error
	}{
		{
			name: "given valid id should return user",
			args: args{id: 1},
			given: func(t *testing.T) error {
				return repo.Create(ctx, newUser(WithEmail("test1@gmail.com")))
			},
			wantErr: nil,
		},
		{
			name: "given invalid id should return error",
			args: args{id: 111},
			given: func(t *testing.T) error {
				return repo.Create(ctx, newUser(WithEmail("test2@gmail.com")))
			},
			wantErr: errors.New("sql: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.given(t)
			assert.NoErr(err)

			u, err := repo.GetById(ctx, tt.args.id)

			assert.Equal(tt.wantErr, err)
			if u != nil {
				assert.Equal(u.ID, tt.args.id)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.NewRelaxed(t)

	// setup test-containers
	tearDown, ctx, bunDb := containers.SetUp(t)
	defer tearDown(t)

	repo := NewRepo(bunDb)

	// setup test cases
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "given valid user should not return error",
			args:    args{u: createUser()},
			wantErr: nil,
		},
		{
			name:    "given nil user should return error",
			args:    args{u: nil},
			wantErr: domain.NilEntityError,
		},
		{
			name:    "given user with id should return error",
			args:    args{u: newUser(WithId(1))},
			wantErr: errors.New("duplicate key value violates unique constraint"),
		},
		{
			name:    "given user with non-unique email should return error",
			args:    args{u: createUser()},
			wantErr: errors.New("duplicate key value violates unique constraint \"users_email_key\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.args.u)

			if tt.wantErr != nil {
				assert.True(strings.Contains(err.Error(), tt.wantErr.Error()))
			} else {
				assert.NoErr(err)
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.NewRelaxed(t)

	// setup test-containers
	tearDown, ctx, bunDb := containers.SetUp(t)
	defer tearDown(t)

	repo := NewRepo(bunDb)

	// setup test cases
	type args struct {
		u *User
	}
	tests := []struct {
		name    string
		args    args
		given   func() error
		verify  func(t *testing.T)
		wantErr error
	}{
		{
			name: "given valid user input should not return error",
			args: args{u: newUser(WithId(1), WithCreatedAt(time.Now()), WithEmail("test@test.com"))},
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
			name: "given nil user input should return error",
			args: args{u: nil},
			given: func() error {
				return nil
			},
			verify:  func(t *testing.T) {},
			wantErr: domain.NilEntityError,
		},
		{
			name: "given user with non existing id should return error",
			args: args{u: newUser(WithId(111), WithCreatedAt(time.Now()))},
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
			err := tt.given()
			assert.NoErr(err)

			err = repo.Update(ctx, tt.args.u)
			if tt.wantErr != nil {
				assert.True(strings.Contains(err.Error(), tt.wantErr.Error()))
			} else {
				assert.NoErr(err)
			}

			tt.verify(t)
		})
	}
}

func newUser(opts ...Option) *User {
	u := createUser()
	for _, opt := range opts {
		opt(u)
	}
	return u
}
