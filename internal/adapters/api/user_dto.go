package api

import (
	"time"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils"
	"github.com/google/uuid"
)

// UserDto represents user DTO.
type UserDto struct {
	ID          string    `json:"id"`
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

func Id(id string) UserDtoOption {
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

func Enabled(e bool) UserDtoOption {
	return func(r *UserDto) {
		r.Enabled = e
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

// ToUser converts UserDto into a User pointer.
func ToUser(r *UserDto) (*user.User, error) {
	var id uuid.UUID
	var err error

	if utils.IsBlank(r.ID) {
		id, err = uuid.NewRandom()
	} else {
		id, err = uuid.Parse(r.ID)
	}

	if err != nil {
		return nil, err
	}

	return user.New(
		user.Id(id),
		user.Email(r.Email),
		user.FullName(r.FullName),
		user.DateOfBirth(r.DateOfBirth),
		user.Location(r.Location),
		user.Enabled(r.Enabled),
		user.Sex(r.Gender.Numberfy()),
	), nil
}

// ToUserDto converts User into a UserDto pointer.
func ToUserDto(u *user.User) *UserDto {
	return NewUserDto(
		Id(u.ID.String()),
		Email(u.Email),
		FullName(u.FullName),
		DateOfBirth(u.DateOfBirth),
		Location(u.Location),
		Sex(u.Gender),
		Enabled(u.Enabled),
	)
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
