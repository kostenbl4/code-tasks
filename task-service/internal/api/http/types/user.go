package types

import (
	"net/http"

	"github.com/kostenbl4/code-tasks/task-service/utils"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetRegisterUserRequest(r *http.Request) (RegisterUserRequest, error) {
	var in RegisterUserRequest

	err := utils.ReadJSON(r, &in)
	if err != nil {
		return RegisterUserRequest{}, err
	}

	// валидация

	return in, nil
}

var GetLoginUserRequest = GetRegisterUserRequest

type LoginUserRequest struct {
	RegisterUserRequest
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
