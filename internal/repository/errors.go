package repository

import "errors"


var (
	// Ошибка, возвращаемая, если задача не найдена
	ErrTaskNotFound = errors.New("task not found")
)
