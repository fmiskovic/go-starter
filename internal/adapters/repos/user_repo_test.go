package repos

import (
	"errors"
	"strings"
	"testing"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"
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
		givenId uuid.UUID
		want    any
		wantErr error
	}{
		{
			name:    "given valid id should return user",
			givenId: uuid.MustParse("220cea28-b2b0-4051-9eb6-9a99e451af01"),
			want:    "john@smith.com",
			wantErr: nil,
		},
		{
			name:    "given non-exisitng id should return error",
			givenId: uuid.MustParse("22222222-b2b0-4051-9eb6-9a99e451af01"),
			want:    nil,
			wantErr: errors.New("sql: no rows in result set"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u, err := repo.GetById(ctx, tt.givenId)

			assert.Equal(tt.wantErr, err)
			if u != nil {
				assert.Equal(u.ID, tt.givenId)
				assert.Equal(u.Email, tt.want)
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
		givenId uuid.UUID
		verify  func(id uuid.UUID, t *testing.T)
		wantErr error
	}{
		{
			name:    "given valid id should delete user",
			givenId: uuid.MustParse("220cea28-b2b0-4051-9eb6-9a99e451af03"),
			verify: func(id uuid.UUID, t *testing.T) {
				u, err := repo.GetById(ctx, id)
				assert.True(strings.Contains(err.Error(), "no rows in result set"))
				assert.True(u == nil)
			},
			wantErr: nil,
		},
		{
			name:    "given non-existing id should not return error",
			givenId: uuid.MustParse("22222222-b2b0-4051-9eb6-9a99e451af01"),
			verify: func(id uuid.UUID, t *testing.T) {
				_, err := repo.GetById(ctx, id)
				assert.True(strings.Contains(err.Error(), "no rows in result set"))
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteById(ctx, tt.givenId)
			assert.Equal(tt.wantErr, err)
			tt.verify(tt.givenId, t)
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
		{
			name: "given valid user with credentials and roles should not return error",
			args: args{u: user.New(
				user.Email("test2@fake.com"),
				user.Credentials(security.NewCredentials(
					security.Username("test2"),
					security.Password("test2"),
				)),
				user.Roles(
					security.NewRole(
						security.RoleName(security.ROLE_USER),
					),
				),
			),
			},
			wantErr: nil,
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
		verify  func(id uuid.UUID, t *testing.T)
		wantErr error
	}{
		{
			name: "given valid user input should not return error",
			args: args{
				u: user.New(user.Email("updated1@fake.com"),
					user.Id(uuid.MustParse("220cea28-b2b0-4051-9eb6-9a99e451af03"))),
			},
			verify: func(id uuid.UUID, t *testing.T) {
				u, err := repo.GetById(ctx, id)
				assert.NoErr(err)
				assert.Equal("updated1@fake.com", u.Email)
			},
			wantErr: nil,
		},
		{
			name:    "given nil user input should return error",
			args:    args{u: nil},
			verify:  func(id uuid.UUID, t *testing.T) {},
			wantErr: ErrNilEntity,
		},
		{
			name: "given user with non-existing id should return error",
			args: args{
				u: user.New(user.Email("updated3@fake.com")),
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
			err := repo.Update(ctx, tt.args.u)
			if tt.wantErr != nil {
				assert.True(strings.Contains(err.Error(), tt.wantErr.Error()))
			} else {
				assert.NoErr(err)
			}

			if tt.args.u != nil {
				tt.verify(tt.args.u.ID, t)
			}
		})
	}
}

func TestUserRepo_GetPage(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	assert := is.NewRelaxed(t)

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
		want    any
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
			want:    "john@smith.com", // value from ./testdata/fixutes.yml
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
			want:    "john@smith.com", // value from ./testdata/fixutes.yml
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestUserRepo_GetByUsername(t *testing.T) {
	// skip in short mode
	if testing.Short() {
		return
	}

	// setup db
	tearDown, ctx, bunDb := testx.SetUpDb(t)
	defer tearDown(t)

	repo := NewUserRepo(bunDb)

	type args struct {
		username string
	}
	tests := []struct {
		name      string
		args      args
		wantEmail any
		wantErr   bool
	}{
		{
			name:      "given valid username should return user",
			args:      args{username: "username1"},
			wantEmail: "john@smith.com",
			wantErr:   false,
		},
		{
			name:      "given non-existing username should return error",
			args:      args{username: "non-existing-username"},
			wantEmail: nil,
			wantErr:   true,
		},
		{
			name:      "given emtpy username should return error",
			args:      args{username: ""},
			wantEmail: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByUsername(ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Email != tt.wantEmail {
				t.Errorf("UserRepo.GetByUsername() = %v, want %v", got, tt.wantEmail)
			}
			if got != nil && got.Credentials != nil && got.Credentials.Username != tt.args.username {
				t.Errorf("UserRepo.GetByUsername() = got username %v, want username %v", got.Credentials.Username, tt.args.username)
			}
		})
	}
}
