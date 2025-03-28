package types

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/utils"
	"errors"
	"net/http"
)

// ErrorResponse - структура для выходных данных в случае ошибки
type ErrorResponse struct {
	Error string `json:"error"`
}

// пока что сделал так более унифицированную обработку ошибок
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
