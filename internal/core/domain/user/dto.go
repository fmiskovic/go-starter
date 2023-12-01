package user

import (
	"time"

	"github.com/fmiskovic/go-starter/internal/core/domain"
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

// GenderDto can be Male, Female and Other.
type GenderDto string

// Numberfy converts GenderDto into a user.Gender.
func (g GenderDto) Numberfy() Gender {
	switch g {
	case "Male":
		return MALE
	case "Female":
		return FEMALE
	case "Other":
		return OTHER
	default:
		return OTHER
	}
}

// ConvertToUser converts User DTO into a User entity.
func ConvertToUser(r *Dto) (*User, error) {
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

	return New(
		Id(id),
		Email(r.Email),
		FullName(r.FullName),
		DateOfBirth(r.DateOfBirth),
		Location(r.Location),
		Enabled(r.Enabled),
		Sex(r.Gender.Numberfy()),
	), nil
}

// ConvertToDto converts User entity into a User DTO.
func ConvertToDto(u *User) *Dto {
	return &Dto{
		ID:          u.ID.String(),
		Email:       u.Email,
		FullName:    u.FullName,
		DateOfBirth: u.DateOfBirth,
		Location:    u.Location,
		Gender:      GenderDto(u.Gender.Stringify()),
		Enabled:     u.Enabled,
	}
}

// ConvertToPageDto converts User entities Page into a User DTO Page.
func ConvertToPageDto(page domain.Page[User]) *domain.Page[Dto] {
	var dtos []Dto
	for _, u := range page.Elements {
		dtos = append(dtos, *ConvertToDto(&u))
	}
	return &domain.Page[Dto]{
		TotalPages:    page.TotalPages,
		TotalElements: page.TotalElements,
		Elements:      dtos,
	}
}
