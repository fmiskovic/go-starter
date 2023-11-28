package security

import (
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:roles,alias:r"`

	domain.Entity
	Name   string    `bun:"name,notnull"`
	UserID uuid.UUID `bun:"user_id,notnull,unique"`
}

func NewRole() *Role {
	r := new(Role)
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
