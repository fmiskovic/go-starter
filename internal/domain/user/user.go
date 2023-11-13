package user

import (
	"github.com/fmiskovic/go-starter/internal/domain"
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	domain.Entity
	Email       string    `bun:"email,notnull,unique" json:"email"`
	FullName    string    `bun:"full_name" json:"fullname"`
	DateOfBirth time.Time `bun:"date_of_birth" json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      Gender    `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

func NewUser(opts ...Option) *User {
	u := &User{Enabled: false}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

type Option func(*User)

func WithEmail(email string) Option {
	return func(u *User) {
		u.Email = email
	}
}

func WithFullName(fullName string) Option {
	return func(u *User) {
		u.FullName = fullName
	}
}

func WithDateOfBirth(dateOfBirth time.Time) Option {
	return func(u *User) {
		u.DateOfBirth = dateOfBirth
	}
}

func WithLocation(location string) Option {
	return func(u *User) {
		u.Location = location
	}
}

func WithGender(gender Gender) Option {
	return func(u *User) {
		u.Gender = gender
	}
}

func WithCreatedAt(time time.Time) Option {
	return func(u *User) {
		u.CreatedAt = time
	}
}

func WithUpdatedAt(time time.Time) Option {
	return func(u *User) {
		u.UpdatedAt = time
	}
}

func WithEnabled(enabled bool) Option {
	return func(u *User) {
		u.Enabled = enabled
	}
}

func (u *User) EnableIt() {
	u.Enabled = true
}

func (u *User) DisableIt() {
	u.Enabled = false
}

type Gender uint8

const (
	MALE Gender = iota
	FEMALE
	OTHER
)
