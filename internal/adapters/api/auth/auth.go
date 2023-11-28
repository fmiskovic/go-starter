package auth

type Request struct {
	Username string `validate:"required,min=3,max=24" json:"username"`
	Password string `validate:"required,min=8,max=72" json:"password"`
}

func NewRequest(username string, password string) Request {
	return Request{Username: username, Password: password}
}

type Response struct {
	Token string `json:"token"`
}

func NewResponse(token string) Response {
	return Response{Token: token}
}
