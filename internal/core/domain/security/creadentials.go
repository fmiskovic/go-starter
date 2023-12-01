package security

import (
	"log/slog"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Credentials struct {
	bun.BaseModel `bun:"table:credentials,alias:c"`

	domain.Entity
	Username string    `bun:"username,notnull,unique"`
	Password string    `bun:"password_hash,notnull"`
	UserID   uuid.UUID `bun:"user_id,notnull,unique"`
}

func NewCredentials(opts ...CredentialsOption) *Credentials {
	// recover in case uuid.New() panic
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Recovered in user.New() when uuid.New() panic", r)
		}
	}()

	crd := &Credentials{Entity: domain.Entity{ID: uuid.New()}}

	for _, opt := range opts {
		opt(crd)
	}
	return crd
}

type CredentialsOption func(crd *Credentials)

func Username(u string) CredentialsOption {
	return func(crd *Credentials) {
		crd.Username = u
	}
}

func Password(p string) CredentialsOption {
	return func(crd *Credentials) {
		crd.Password = p
	}
}
