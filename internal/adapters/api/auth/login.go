package auth

type LoginRequest struct {
	Username string
	Password string
}

func NewLoginRequest(username string, password string) LoginRequest {
	return LoginRequest{Username: username, Password: password}
}
