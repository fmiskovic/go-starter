package user

import (
	"context"
	"github.com/fmiskovic/go-starter/domain"

	"github.com/fmiskovic/go-starter/database"
	"github.com/uptrace/bun"
)

// Repo is user implementation of core.Repo interface
type Repo struct {
	db *bun.DB
}

func NewRepo(db *bun.DB) *Repo {
	return &Repo{db}
}

func NewDefaultRepo() *Repo {
	return &Repo{database.ConnectDefault()}
}

// GetById returns user queried by id
func (repo *Repo) GetById(ctx context.Context, id int64) (*User, error) {
	var user = &User{}

	err := repo.db.NewSelect().Model(user).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Save persists new user or updates existing one
func (repo *Repo) Save(ctx context.Context, u *User) error {
	if u.ID > 0 {
		// update
		if _, err := repo.db.NewUpdate().Model(u).Where("id = ?", u.ID).Exec(ctx); err != nil {
			return err
		}

		return nil
	}

	// save
	if _, err := repo.db.NewInsert().Model(u).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteById deletes user by id
func (repo *Repo) DeleteById(ctx context.Context, id int64) error {
	if _, err := repo.db.NewDelete().Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// GetPage of users
func (repo *Repo) GetPage(ctx context.Context, p domain.Pageable) (domain.Page[User], error) {
	// TODO
	return domain.Page[User]{}, nil
}
