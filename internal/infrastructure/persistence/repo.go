// Package persistence contains repository implementations.
package persistence

import (
	"context"
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/fmiskovic/go-starter/internal/interfaces/pagination"
)

// Repo is generic repository.
type Repo[ID any, T any] interface {
	GetById(ctx context.Context, id ID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	DeleteById(ctx context.Context, id ID) error
	GetPage(ctx context.Context, p pagination.Pageable) (domain.Page[T], error)
}
