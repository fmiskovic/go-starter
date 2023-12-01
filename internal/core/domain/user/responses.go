package user

type SignInResponse struct {
	Token string `json:"token"`
}

type SignUpResponse struct {
	ID string `json:"id"`
}

type CreateResponse struct {
	Dto
}

type UpdateResponse struct {
	Dto
}
