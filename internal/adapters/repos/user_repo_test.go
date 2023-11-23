package repos

import (
	"errors"
	"strings"
	"testing"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestUserRepo_GetById(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.New(t)

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	// setup test cases
	tests := []struct {
		name    string
		given   func(t *testing.T) (uuid.UUID, error)
		wantErr error
	}{
		{
			name: "given valid id should return user",
			given: func(t *testing.T) (uuid.UUID, error) {
				u := user.New(user.Email("test1@gmail.com"))
				err := repo.Create(ctx, u)
				return u.ID, err
			},
			wantErr: nil,
		},
		{
			name: "given invalid id should return error",
			given: func(t *testing.T) (uuid.UUID, error) {
				err := repo.Create(ctx, user.New(user.Email("test2@gmail.com")))
				return uuid.New(), err
			},
			wantErr: errors.New("sql: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tt.given(t)
			assert.NoErr(err)

			u, err := repo.GetById(ctx, id)

			assert.Equal(tt.wantErr, err)
			if u != nil {
				assert.Equal(u.ID, id)
			}
		})
	}
}

func TestUserRepo_DeleteById(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.New(t)

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	// setup test cases
	tests := []struct {
		name    string
		given   func(t *testing.T) (uuid.UUID, error)
		verify  func(id uuid.UUID, t *testing.T)
		wantErr error
	}{
		{
			name: "given valid id should delete user",
			given: func(t *testing.T) (uuid.UUID, error) {
				u := user.New(user.Email("test1@gmail.com"))
				err := repo.Create(ctx, u)
				return u.ID, err
			},
			verify: func(id uuid.UUID, t *testing.T) {
				u, err := repo.GetById(ctx, id)
				assert.True(strings.Contains(err.Error(), "no rows in result set"))
				assert.True(u == nil)
			},
			wantErr: nil,
		},
		{
			name: "given invalid id should not return error",
			given: func(t *testing.T) (uuid.UUID, error) {
				err := repo.Create(ctx, user.New(user.Email("test2@gmail.com")))
				return uuid.New(), err
			},
			verify: func(id uuid.UUID, t *testing.T) {
				_, err := repo.GetById(ctx, id)
				assert.True(strings.Contains(err.Error(), "no rows in result set"))
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tt.given(t)
			assert.NoErr(err)

			err = repo.DeleteById(ctx, id)
			assert.Equal(tt.wantErr, err)
			tt.verify(id, t)
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.New(t)

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	// setup test cases
	type args struct {
		u *user.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "given valid user should not return error",
			args:    args{u: user.New(user.Email("test1@fake.com"))},
			wantErr: nil,
		},
		{
			name:    "given nil user should return error",
			args:    args{u: nil},
			wantErr: ErrNilEntity,
		},
		{
			name:    "given user with non-unique email should return error",
			args:    args{u: user.New(user.Email("test1@fake.com"))},
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

	assert := is.New(t)

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	// setup test cases
	type args struct {
		u *user.User // for update
	}
	tests := []struct {
		name    string
		args    args
		given   func() (*user.User, error)
		verify  func(id uuid.UUID, t *testing.T)
		wantErr error
	}{
		{
			name: "given valid user input should not return error",
			args: args{
				u: user.New(user.Email("updated1@fake.com")),
			},
			given: func() (*user.User, error) {
				u := user.New(user.Email("test1@fake.com"))
				err := repo.Create(ctx, u)
				return u, err
			},
			verify: func(id uuid.UUID, t *testing.T) {
				u, err := repo.GetById(ctx, id)
				assert.NoErr(err)
				assert.Equal("updated1@fake.com", u.Email)
			},
			wantErr: nil,
		},
		{
			name: "given nil user input should return error",
			args: args{u: nil},
			given: func() (*user.User, error) {
				return user.New(), nil
			},
			verify:  func(id uuid.UUID, t *testing.T) {},
			wantErr: ErrNilEntity,
		},
		{
			name: "given user with non existing id should return error",
			args: args{
				u: user.New(user.Email("updated3@fake.com")),
			},
			given: func() (*user.User, error) {
				return user.New(), nil
			},
			verify: func(id uuid.UUID, t *testing.T) {
				_, err := repo.GetById(ctx, id)
				assert.True(strings.Contains(err.Error(), "no rows in result set"))
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := tt.given()
			assert.NoErr(err)

			if tt.args.u != nil && u != nil {
				tt.args.u.ID = u.ID
			}

			err = repo.Update(ctx, tt.args.u)
			if tt.wantErr != nil {
				assert.True(strings.Contains(err.Error(), tt.wantErr.Error()))
			} else {
				assert.NoErr(err)
			}

			tt.verify(u.ID, t)
		})
	}
}

func TestUserRepo_GetPage(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.New(t)

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	// setup test cases
	type args struct {
		pageable domain.Pageable
	}
	tests := []struct {
		name    string
		args    args
		given   func(t *testing.T) error
		want    string
		wantErr error
	}{
		{
			name: "given page request should return page of users",
			args: args{
				pageable: domain.Pageable{
					Offset: 0,
					Size:   5,
					Sort:   domain.NewSort(domain.NewOrder(domain.WithProperty("email"))),
				},
			},
			given: func(t *testing.T) error {
				u := user.New(user.Email("test11@gmail.com"))
				return repo.Create(ctx, u)
			},
			want:    "test11@gmail.com",
			wantErr: nil,
		},
		{
			name: "given page request without sort should return page of users",
			args: args{
				pageable: domain.Pageable{
					Offset: 0,
					Size:   5,
				},
			},
			given: func(t *testing.T) error {
				u := user.New(user.Email("test12@gmail.com"))
				return repo.Create(ctx, u)
			},
			want:    "test11@gmail.com",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.given(t)
			assert.NoErr(err)

			p, err := repo.GetPage(ctx, tt.args.pageable)

			assert.Equal(tt.wantErr, err)
			if err == nil {
				assert.True(len(p.Elements) > 0)
				assert.True(p.TotalPages == 1)
				assert.Equal(p.Elements[0].Email, tt.want)
			}
		})
	}
}
