package user

import (
	"github.com/fmiskovic/go-starter/internal/domain"
	"time"

	"github.com/uptrace/bun"
)

// User represents database entity.
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	domain.Entity
	Email       string    `bun:"email,notnull,unique"`
	FullName    string    `bun:"full_name,nullzero"`
	DateOfBirth time.Time `bun:"date_of_birth,nullzero"`
	Location    string    `bun:"location,nullzero"`
	Gender      Gender    `bun:"gender,nullzero"`
	Enabled     bool      `bun:"enabled"`
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
