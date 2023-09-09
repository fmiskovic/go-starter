package data

import "time"

type UserProfile struct {
	ID          int64     `bun:",pk,autoincrement" json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      Gender    `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

type Gender uint8

const (
	MALE Gender = iota
	FEMALE
	OTHER
)
