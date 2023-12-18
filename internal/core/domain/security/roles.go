package security

import (
	"log/slog"
	"time"

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

func NewRole(name string) *Role {
	// recover in case uuid.New() panic
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Recovered in user.New() when uuid.New() panic", r)
		}
	}()

	now := time.Now()

	return &Role{
		Entity: domain.Entity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name: name,
	}
}

type RoleOption func(*Role)

func RoleName(name string) RoleOption {
	return func(r *Role) {
		r.Name = name
	}
}
