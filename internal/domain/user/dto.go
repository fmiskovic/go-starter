package user

import (
	"github.com/fmiskovic/go-starter/internal/domain"
	"time"
)

// Dto represents user create and update request data
type Dto struct {
	ID          uint64    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Email       string    `validate:"required,min=3" json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      Gender    `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

// NewDto constructs new user request
func NewDto(opts ...DtoOption) *Dto {
	r := &Dto{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type DtoOption func(req *Dto)

func Id(id uint64) DtoOption {
	return func(r *Dto) {
		r.ID = id
	}
}

func Email(email string) DtoOption {
	return func(r *Dto) {
		r.Email = email
	}
}

func FullName(fullName string) DtoOption {
	return func(r *Dto) {
		r.FullName = fullName
	}
}

func DateOfBirth(dateOfBirth time.Time) DtoOption {
	return func(r *Dto) {
		r.DateOfBirth = dateOfBirth
	}
}

func Location(location string) DtoOption {
	return func(r *Dto) {
		r.Location = location
	}
}

func Sex(gender Gender) DtoOption {
	return func(r *Dto) {
		r.Gender = gender
	}
}

func toUser(r *Dto) *User {
	return &User{
		Entity: domain.Entity{
			ID: r.ID,
		},
		Email:       r.Email,
		FullName:    r.FullName,
		DateOfBirth: r.DateOfBirth,
		Location:    r.Location,
		Gender:      r.Gender,
		Enabled:     r.Enabled,
	}
}

func toDto(u *User) *Dto {
	return &Dto{
		ID:          u.ID,
		Email:       u.Email,
		FullName:    u.FullName,
		DateOfBirth: u.DateOfBirth,
		Location:    u.Location,
		Gender:      u.Gender,
		Enabled:     u.Enabled,
	}
}
