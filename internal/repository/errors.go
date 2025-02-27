package repository

import "errors"

var (
	// Ошибка, возвращаемая, если задача не найдена
	ErrTaskNotFound = errors.New("task not found")
	// Ошибка, возвращаемая, если пользователь не найден
	ErrUserNotFound = errors.New("user not found")
	// Ошибка, возвращаемая, если сессия не найдена
	ErrSessionNotFound = errors.New("session not found")
)
