package user

import (
	"github.com/fmiskovic/go-starter/internal/domain"
	"time"
)

// Request represents user create and update request data
type Request struct {
	ID          uint64    `json:"id"`
	Email       string    `validate:"required,min=3" json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      Gender    `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

// NewRequest constructs new user request
func NewRequest(opts ...RequestOption) *Request {
	r := &Request{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type RequestOption func(req *Request)

func Id(id uint64) RequestOption {
	return func(r *Request) {
		r.ID = id
	}
}

func Email(email string) RequestOption {
	return func(r *Request) {
		r.Email = email
	}
}

func FullName(fullName string) RequestOption {
	return func(r *Request) {
		r.FullName = fullName
	}
}

func DateOfBirth(dateOfBirth time.Time) RequestOption {
	return func(r *Request) {
		r.DateOfBirth = dateOfBirth
	}
}

func Location(location string) RequestOption {
	return func(r *Request) {
		r.Location = location
	}
}

func Sex(gender Gender) RequestOption {
	return func(r *Request) {
		r.Gender = gender
	}
}

func toUser(r *Request) *User {
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
