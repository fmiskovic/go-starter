package data

import "time"

type User struct {
	ID          int64     `bun:",pk,autoincrement" json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      Gender    `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

func NewUser(opts ...UserOption) *User {
	u := &User{Enabled: false}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

type UserOption func(*User)

func WithEmail(email string) UserOption {
	return func(u *User) {
		u.Email = email
	}
}

func WithFullName(fullName string) UserOption {
	return func(u *User) {
		u.FullName = fullName
	}
}

func WithDateOfBirth(dateOfBirth time.Time) UserOption {
	return func(u *User) {
		u.DateOfBirth = dateOfBirth
	}
}

func WithLocation(location string) UserOption {
	return func(u *User) {
		u.Location = location
	}
}

func WithGender(gender Gender) UserOption {
	return func(u *User) {
		u.Gender = gender
	}
}

func WithEnabled(enabled bool) UserOption {
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
