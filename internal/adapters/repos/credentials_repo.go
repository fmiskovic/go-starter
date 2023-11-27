package repos

import (
	"context"
	"time"

	"github.com/fmiskovic/go-starter/internal/core/domain/security"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// CredentialsRepo is implementation of ports.CredentialsRepo.
type CredentialsRepo struct {
	db *bun.DB
}

// NewCredentialsRepo instantiate new CredentialsRepo.
func NewCredentialsRepo(db *bun.DB) CredentialsRepo {
	return CredentialsRepo{db}
}

// GetById returns credentials by specified id.
func (repo CredentialsRepo) GetById(ctx context.Context, id uuid.UUID) (*security.Credentials, error) {
	var c = &security.Credentials{}

	err := repo.db.NewSelect().Model(c).Where("? = ?", bun.Ident("id"), id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Create persists new credentials entity.
func (repo CredentialsRepo) Create(ctx context.Context, c *security.Credentials) error {
	if c == nil {
		return ErrNilEntity
	}

	if _, err := repo.db.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// Update existing credentials entity.
func (repo CredentialsRepo) Update(ctx context.Context, c *security.Credentials) error {
	if c == nil {
		return ErrNilEntity
	}
	c.UpdatedAt = time.Now()
	if _, err := repo.db.NewUpdate().Model(c).OmitZero().Where("id = ?", c.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteById remove credentials entity by specified id.
func (repo CredentialsRepo) DeleteById(ctx context.Context, id uuid.UUID) error {
	if _, err := repo.db.NewDelete().Model(new(security.Credentials)).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}
