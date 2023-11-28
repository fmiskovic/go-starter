package ports

import (
	"context"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
)

// BaseRepo is generic repository.
type BaseRepo[ID any, T any] interface {
	GetById(ctx context.Context, id ID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	DeleteById(ctx context.Context, id ID) error
}

// UserRepo represents user repository interface.
type UserRepo[ID any] interface {
	GetById(ctx context.Context, id ID) (*user.User, error)
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	DeleteById(ctx context.Context, id ID) error
	GetPage(ctx context.Context, p domain.Pageable) (domain.Page[user.User], error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
}

// CredentialsRepo represents credentials repository interface.
type CredentialsRepo[ID any] interface {
	GetById(ctx context.Context, id ID) (*security.Credentials, error)
	Create(ctx context.Context, crd *security.Credentials) error
	Update(ctx context.Context, crd *security.Credentials) error
	DeleteById(ctx context.Context, id ID) error
}

// RoleRepo represents role repository interface.
type RoleRepo[ID any] interface {
	GetById(ctx context.Context, id ID) (*security.Role, error)
	Create(ctx context.Context, crd *security.Role) error
	Update(ctx context.Context, crd *security.Role) error
	DeleteById(ctx context.Context, id ID) error
}
