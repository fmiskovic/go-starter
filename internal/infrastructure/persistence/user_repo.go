package persistence

import (
	"context"
	"github.com/fmiskovic/go-starter/internal/domain"
	"github.com/fmiskovic/go-starter/internal/domain/user"
	"github.com/fmiskovic/go-starter/internal/interfaces/api"
	"github.com/fmiskovic/go-starter/internal/interfaces/pagination"
	"time"

	"github.com/uptrace/bun"
)

// UserRepo is implementation of generic Repo interface.
type UserRepo struct {
	db *bun.DB
}

// NewUserRepo instantiate new UserRepo.
func NewUserRepo(db *bun.DB) UserRepo {
	return UserRepo{db}
}

// GetById returns user by specified id.
func (repo *UserRepo) GetById(ctx context.Context, id uint64) (*user.User, error) {
	var u = &user.User{}

	err := repo.db.NewSelect().Model(u).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Create persists new user entity.
func (repo *UserRepo) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return api.NilUserError
	}
	if _, err := repo.db.NewInsert().Model(u).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Update existing persisted user entity.
func (repo *UserRepo) Update(ctx context.Context, u *user.User) error {
	if u == nil {
		return api.NilUserError
	}
	u.UpdatedAt = time.Now()
	if _, err := repo.db.NewUpdate().Model(u).OmitZero().Where("id = ?", u.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteById remove user entity by specified id.
func (repo *UserRepo) DeleteById(ctx context.Context, id uint64) error {
	if _, err := repo.db.NewDelete().Model(new(user.User)).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// GetPage respond with a page of users.
func (repo *UserRepo) GetPage(ctx context.Context, p pagination.Pageable) (domain.Page[user.User], error) {
	var users []user.User
	count, err := repo.db.
		NewSelect().
		Model(&users).
		Limit(p.Size).
		Offset(p.Offset).
		Order(pagination.StringifyOrders(p.Sort)...).
		ScanAndCount(ctx)

	totalPages := 0
	if count != 0 && p.Size != 0 {
		totalPages = (len(users) + p.Size - 1) / p.Size
	}

	return domain.Page[user.User]{
		TotalPages:    totalPages,
		TotalElements: count,
		Elements:      users,
	}, err
}
