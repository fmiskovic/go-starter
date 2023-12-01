package user

import "time"

const (
	//TimeFormat1 to format date into
	TimeFormat1 = "2000-01-01"
	//TimeFormat2 Other format to format date time
	TimeFormat2 = "January 01, 2000"
)

type Request struct {
	Email       string    `validate:"required,min=3" json:"email"`
	FullName    string    `json:"fullname"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Location    string    `json:"location"`
	Gender      GenderDto `json:"gender"`
}

type SignInRequest struct {
	Username string `validate:"required,min=3,max=24" json:"username"`
	Password string `validate:"required,min=8,max=72" json:"password"`
}

type ChangePasswordRequest struct {
	ID       string `json:"id"`
	Password string `validate:"required,min=8,max=72" json:"password"`
}

type ConfirmEmailRequest struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type CreateRequest struct {
	Username string `validate:"required,min=3,max=24" json:"username"`
	Password string `validate:"required,min=8,max=72" json:"password"`
	Request
}

type UpdateRequest struct {
	ID string `json:"id"`
	Request
}
