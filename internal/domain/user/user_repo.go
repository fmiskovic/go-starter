package user

import (
	"context"
	"github.com/fmiskovic/go-starter/internal/domain"

	"github.com/uptrace/bun"
)

// UserRepo is implementation of core.UserRepo interface
type UserRepo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) UserRepo {
	return UserRepo{db}
}

// GetById returns user by id
func (repo *UserRepo) GetById(ctx context.Context, id int64) (*User, error) {
	var user = &User{}

	err := repo.db.NewSelect().Model(user).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Create new user
func (repo *UserRepo) Create(ctx context.Context, u *User) error {
	if u == nil {
		return domain.NilEntityError
	}
	if _, err := repo.db.NewInsert().Model(u).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Update existing user
func (repo *UserRepo) Update(ctx context.Context, u *User) error {
	if u == nil {
		return domain.NilEntityError
	}
	if _, err := repo.db.NewUpdate().Model(u).Where("id = ?", u.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteById remove user by id
func (repo *UserRepo) DeleteById(ctx context.Context, id int64) error {
	if _, err := repo.db.NewDelete().Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// GetPage of users
func (repo *UserRepo) GetPage(ctx context.Context, p domain.Pageable) (domain.Page[User], error) {
	var users []User
	err := repo.db.
		NewSelect().
		Model(&User{}).
		Limit(p.Size).
		Offset(p.Offset).
		Order(domain.Orders(p.Sort)...).
		Scan(ctx, users)

	totalCount := 0
	totalPages := 0
	if p.IncludeTotalCount {
		totalCount, err = repo.db.
			NewSelect().
			Model(&User{}).
			Count(ctx)

		if totalCount != 0 && p.Offset != 0 {
			totalPages = totalCount / p.Offset
		}
	}

	return domain.Page[User]{
		TotalPages:    totalPages,
		TotalElements: totalCount,
		Elements:      users,
	}, err
}
