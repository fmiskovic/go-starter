package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apiErr "github.com/fmiskovic/go-starter/internal/core/error"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils/password"
	"github.com/google/uuid"

	"github.com/uptrace/bun"
)

// UserRepo is implementation of ports.UserRepo interface.
type UserRepo struct {
	db *bun.DB
}

// NewUserRepo instantiate new UserRepo.
func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{db}
}

// GetById returns user by specified id.
func (repo *UserRepo) GetById(ctx context.Context, id uuid.UUID) (*user.User, error) {
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
func (repo *UserRepo) Update(ctx context.Context, u *user.User) error {
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
func (repo *UserRepo) DeleteById(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.db.NewDelete().Model(new(user.User)).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// GetPage respond with a page of users.
func (repo *UserRepo) GetPage(ctx context.Context, p domain.Pageable) (domain.Page[user.User], error) {
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

// GetByUsername returns user by username.
func (repo *UserRepo) GetByUsername(ctx context.Context, username string) (*user.User, error) {
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

// ChangePassword updates users password.
func (repo *UserRepo) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest) error {
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// repo.mutex.Lock()
		// defer repo.mutex.Unlock()

		var u = new(user.User)

		err := tx.NewSelect().
			Model(u).
			Relation("Credentials", func(sq *bun.SelectQuery) *bun.SelectQuery {
				return sq.Where("username = ?", req.Username)
			}).
			Scan(ctx)

		if err != nil {
			return err
		}

		if !password.CheckPasswordHash(req.OldPassword, u.Credentials.Password) {
			return errors.New("invalid old password")
		}

		newPwd, err := password.HashPassword(req.NewPassword)
		if err != nil {
			return err
		}

		crd := u.Credentials
		crd.Password = newPwd
		crd.UpdatedAt = time.Now()
		u.UpdatedAt = crd.UpdatedAt

		if _, err := tx.NewUpdate().Model(crd).OmitZero().Where("user_id = ?", u.ID).Exec(ctx); err != nil {
			return err
		}

		return nil
	})
}

// AddRoles to existing user.
func (repo *UserRepo) AddRoles(ctx context.Context, roleNames []string, id uuid.UUID) error {
	l := len(roleNames)
	if l == 0 {
		return nil
	}
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		exists, err := tx.NewSelect().
			Model((*user.User)(nil)).
			Where("? = ?", bun.Ident("id"), id).
			Exists(ctx)

		if err != nil {
			return err
		}
		if !exists {
			return apiErr.ErrInvalidId
		}

		var roles = make([]*security.Role, l)
		for i, name := range roleNames {
			role := security.NewRole(name)
			role.UserID = id
			roles[i] = role
		}

		if len(roles) > 0 {
			if _, err := tx.NewInsert().Model(&roles).Exec(ctx); err != nil {
				return err
			}

			u := &user.User{Entity: domain.Entity{ID: id, UpdatedAt: time.Now()}}
			if _, err := tx.NewUpdate().Model((u)).OmitZero().Where("? = ?", bun.Ident("id"), id).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

// RemoveRoles from existing user.
func (repo *UserRepo) RemoveRoles(ctx context.Context, roleNames []string, id uuid.UUID) error {
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		exists, err := tx.NewSelect().
			Model((*user.User)(nil)).
			Where("? = ?", bun.Ident("id"), id).
			Exists(ctx)

		if err != nil {
			return err
		}
		if !exists {
			return apiErr.ErrInvalidId
		}

		if len(roleNames) > 0 {
			_, err = tx.NewDelete().
				Model(&security.Role{}).
				Where("user_id = ?", id).
				Where("name IN (?)", bun.In(roleNames)).
				Exec(ctx)

			if err != nil {
				return err
			}

			u := &user.User{Entity: domain.Entity{ID: id, UpdatedAt: time.Now()}}
			if _, err := tx.NewUpdate().Model((u)).OmitZero().Where("? = ?", bun.Ident("id"), id).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
}

// EnableDisable enables user if it is disabled or vice versa.
func (repo *UserRepo) EnableDisable(ctx context.Context, id uuid.UUID) error {
	return repo.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// repo.mutex.Lock()
		// defer repo.mutex.Unlock()

		var u = &user.User{}

		if err := tx.NewSelect().Model(u).Column("enabled").Where("? = ?", bun.Ident("id"), id).Scan(ctx); err != nil {
			return err
		}

		u.UpdatedAt = time.Now()
		u.Enabled = !u.Enabled

		if _, err := tx.NewUpdate().Model(u).Column("enabled", "updated_at").Where("id = ?", id).Exec(ctx); err != nil {
			return err
		}
		return nil
	})
}
