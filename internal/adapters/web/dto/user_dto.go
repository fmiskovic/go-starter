// Package dto defines API contracts using DTOs, request, and response structs.
package dto

import (
	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"time"
)

// UserDto represents user DTO.
type UserDto struct {
	ID          uint64    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Email       string    `validate:"required,min=3" json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      GenderDto `json:"gender"`
	Enabled     bool      `json:"enabled"`
}

// NewUserDto instantiate new user DTO.
func NewUserDto(opts ...UserDtoOption) *UserDto {
	r := &UserDto{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type UserDtoOption func(req *UserDto)

func Id(id uint64) UserDtoOption {
	return func(r *UserDto) {
		r.ID = id
	}
}

func Email(email string) UserDtoOption {
	return func(r *UserDto) {
		r.Email = email
	}
}

func FullName(fullName string) UserDtoOption {
	return func(r *UserDto) {
		r.FullName = fullName
	}
}

func DateOfBirth(dateOfBirth time.Time) UserDtoOption {
	return func(r *UserDto) {
		r.DateOfBirth = dateOfBirth
	}
}

func Location(location string) UserDtoOption {
	return func(r *UserDto) {
		r.Location = location
	}
}

func Sex(g user.Gender) UserDtoOption {
	return func(r *UserDto) {
		r.Gender = GenderDto(g.Stringify())
	}
}

// GenderDto can be Male, Female and Other.
type GenderDto string

// Numberfy converts GenderDto into a user.Gender.
func (g GenderDto) Numberfy() user.Gender {
	switch g {
	case "Male":
		return user.MALE
	case "Female":
		return user.FEMALE
	case "Other":
		return user.OTHER
	default:
		return user.OTHER
	}
}

// ToUser converts UserDto into a user.User pointer.
func ToUser(r *UserDto) *user.User {
	return &user.User{
		Entity: domain.Entity{
			ID: r.ID,
		},
		Email:       r.Email,
		FullName:    r.FullName,
		DateOfBirth: r.DateOfBirth,
		Location:    r.Location,
		Gender:      r.Gender.Numberfy(),
		Enabled:     r.Enabled,
	}
}

// ToUser converts User into a UserDto pointer.
func ToUserDto(u *user.User) *UserDto {
	return &UserDto{
		ID:          u.ID,
		Email:       u.Email,
		FullName:    u.FullName,
		DateOfBirth: u.DateOfBirth,
		Location:    u.Location,
		Gender:      GenderDto(u.Gender.Stringify()),
		Enabled:     u.Enabled,
	}
}

// ToPageDto converts internal Page into a DtoPage.
func ToPageDto(page domain.Page[user.User]) domain.Page[UserDto] {
	var dtos []UserDto
	for _, u := range page.Elements {
		dtos = append(dtos, *ToUserDto(&u))
	}
	return domain.Page[UserDto]{
		TotalPages:    page.TotalPages,
		TotalElements: page.TotalElements,
		Elements:      dtos,
	}
}
