package repos

import (
	"context"
	"database/sql"
	"time"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/google/uuid"

	"github.com/uptrace/bun"
)

// UserRepo is implementation of ports.UserRepo interface.
type UserRepo struct {
	db *bun.DB
}

// NewUserRepo instantiate new UserRepo.
func NewUserRepo(db *bun.DB) UserRepo {
	return UserRepo{db}
}

// GetById returns user by specified id.
func (repo UserRepo) GetById(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var u = &user.User{}

	err := repo.db.NewSelect().Model(u).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Create persists new user entity.
func (repo UserRepo) Create(ctx context.Context, u *user.User) error {
	if u == nil {
		return ErrNilEntity
	}

	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(u).Exec(ctx)
		if err != nil {
			return err
		}
		if u.Credentials != nil {
			_, err = tx.NewInsert().Model(u.Credentials).Exec(ctx)
			if err != nil {
				return err
			}
		}
		if u.Roles != nil {
			_, err = tx.NewInsert().Model(&u.Roles).Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Update existing persisted user entity.
func (repo UserRepo) Update(ctx context.Context, u *user.User) error {
	if u == nil {
		return ErrNilEntity
	}
	u.UpdatedAt = time.Now()
	if _, err := repo.db.NewUpdate().Model(u).OmitZero().Where("id = ?", u.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteById remove user entity by specified id.
func (repo UserRepo) DeleteById(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.db.NewDelete().Model(new(user.User)).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// GetPage respond with a page of users.
func (repo UserRepo) GetPage(ctx context.Context, p domain.Pageable) (domain.Page[user.User], error) {
	var users []user.User
	count, err := repo.db.
		NewSelect().
		Model(&users).
		Limit(p.Size).
		Offset(p.Offset).
		Order(domain.StringifyOrders(p.Sort)...).
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

func (repo UserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u = new(user.User)

	err := repo.db.NewSelect().
		Model(u).
		Relation("Roles").
		Relation("Credentials", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("username = ?", username)
		}).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return u, nil
}
