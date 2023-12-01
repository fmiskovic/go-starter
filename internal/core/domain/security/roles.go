package security

import (
	"log/slog"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var (
	ROLE_ADMIN = "ROLE_ADMIN"
	ROLE_USER  = "ROLE_USER"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`

	domain.Entity
	Name   string    `bun:"name,notnull"`
	UserID uuid.UUID `bun:"user_id,notnull,unique"`
}

func NewRole(opts ...RoleOption) *Role {
	// recover in case uuid.New() panic
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Recovered in user.New() when uuid.New() panic", r)
		}
	}()
	r := &Role{Entity: domain.Entity{ID: uuid.New()}}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type RoleOption func(*Role)

func Id(id uuid.UUID) RoleOption {
	return func(r *Role) {
		r.ID = id
	}
}

func Name(name string) RoleOption {
	return func(r *Role) {
		r.Name = name
	}
}

func RoleUserId(userId uuid.UUID) RoleOption {
	return func(r *Role) {
		r.UserID = userId
	}
}
