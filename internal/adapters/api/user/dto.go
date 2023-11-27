package user

import (
	"time"

	"github.com/fmiskovic/go-starter/internal/core/domain"
	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils"
	"github.com/google/uuid"
)

// UserDto represents user DTO.
type Dto struct {
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

// NewDto instantiate new user DTO.
func NewDto(opts ...DtoOption) *Dto {
	r := &Dto{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

type DtoOption func(req *Dto)

func Id(id string) DtoOption {
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

func Sex(g user.Gender) DtoOption {
	return func(r *Dto) {
		r.Gender = GenderDto(g.Stringify())
	}
}

func Enabled(e bool) DtoOption {
	return func(r *Dto) {
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
func ToUser(r *Dto) (*user.User, error) {
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

// ToDto converts User into a UserDto pointer.
func ToDto(u *user.User) *Dto {
	return NewDto(
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
func ToPageDto(page domain.Page[user.User]) domain.Page[Dto] {
	var dtos []Dto
	for _, u := range page.Elements {
		dtos = append(dtos, *ToDto(&u))
	}
	return domain.Page[Dto]{
		TotalPages:    page.TotalPages,
		TotalElements: page.TotalElements,
		Elements:      dtos,
	}
}
