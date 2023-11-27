package user

import (
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/security"

	"github.com/uptrace/bun"
)

// User represents database entity.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	domain.Entity
	Email       string                `bun:"email,notnull,unique"`
	FullName    string                `bun:"full_name,nullzero"`
	DateOfBirth time.Time             `bun:"date_of_birth,nullzero"`
	Location    string                `bun:"location,nullzero"`
	Gender      Gender                `bun:"gender,nullzero"`
	Enabled     bool                  `bun:"enabled"`
	Credentials *security.Credentials `bun:"rel:has-one,join:id=user_id"`
}

func New(opts ...Option) *User {
	// recover in case uuid.New() panic
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("Recovered in user.New() when uuid.New() panic", r)
		}
	}()
	id := uuid.New()
	u := &User{Entity: domain.Entity{ID: id}}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

type Option func(*User)

func Id(id uuid.UUID) Option {
	return func(u *User) {
		u.ID = id
	}
}

func Email(e string) Option {
	return func(u *User) {
		u.Email = e
	}
}

func FullName(fn string) Option {
	return func(u *User) {
		u.FullName = fn
	}
}

func Location(l string) Option {
	return func(u *User) {
		u.Location = l
	}
}

func Enabled(e bool) Option {
	return func(u *User) {
		u.Enabled = e
	}
}

func DateOfBirth(dob time.Time) Option {
	return func(u *User) {
		u.DateOfBirth = dob
	}
}

func Sex(g Gender) Option {
	return func(u *User) {
		u.Gender = g
	}
}

func Credentials(crd *security.Credentials) Option {
	return func(u *User) {
		u.Credentials = crd
	}
}

// Gender is either MALE, FEMALE or OTHER.
type Gender uint8

// Stringify converts Gender into string.
func (g Gender) Stringify() string {
	switch g {
	case 0:
		return "Male"
	case 1:
		return "Female"
	case 2:
		return "Other"
	default:
		return "Other"
	}
}

const (
	MALE Gender = iota
	FEMALE
	OTHER
)
