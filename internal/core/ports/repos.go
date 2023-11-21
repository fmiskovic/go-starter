package ports

import (
	"context"
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
)

// BaseRepo is generic repository.
type BaseRepo[ID any, T any] interface {
	GetById(ctx context.Context, id ID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	DeleteById(ctx context.Context, id ID) error
	GetPage(ctx context.Context, p Pageable) (domain.Page[T], error)
}

// UserRepo represents user repository interface.
type UserRepo[ID any] interface {
	GetById(ctx context.Context, id ID) (*user.User, error)
	Create(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
	DeleteById(ctx context.Context, id ID) error
	GetPage(ctx context.Context, p Pageable) (domain.Page[user.User], error)
}
