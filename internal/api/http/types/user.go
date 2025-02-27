package types

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	RegisterUserRequest
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
