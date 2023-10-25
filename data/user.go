package data

import "time"

type User struct {
	Entity
	Email       string    `bun:"email,unique" json:"email"`
	FullName    string    `bun:"@full_name" json:"fullname"`
	DateOfBirth time.Time `bun:"@date_of_birth" json:"dateOfBirth"`
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

func WithCreatedAt(time time.Time) UserOption {
	return func(u *User) {
		u.CreatedAt = time
	}
}

func WithUpdatedAt(time time.Time) UserOption {
	return func(u *User) {
		u.UpdatedAt = time
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
