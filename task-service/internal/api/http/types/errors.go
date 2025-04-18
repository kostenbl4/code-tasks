package types

import (
	"errors"
	"net/http"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/utils"
)

// ErrorResponse - структура для выходных данных в случае ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleError(w http.ResponseWriter, err error) {

	if err == nil {
		return
	}

	switch {
	case errors.Is(err, domain.ErrBadRequest),
		errors.Is(err, domain.ErrSessionNotFound),
		errors.Is(err, domain.ErrUserAlreadyExists):
		utils.WriteJSON(w, ErrorResponse{Error: err.Error()}, http.StatusBadRequest)
	case errors.Is(err, domain.ErrIncorrectPassword),
		errors.Is(err, domain.ErrUserNotFound):
		utils.WriteJSON(w, ErrorResponse{Error: "unauthorized"}, http.StatusUnauthorized)
	case errors.Is(err, domain.ErrTaskNotFound):
		utils.WriteJSON(w, ErrorResponse{Error: err.Error()}, http.StatusNotFound)
	default:
		utils.WriteJSON(w, ErrorResponse{Error: "unknown error"}, http.StatusInternalServerError)
	}
}
