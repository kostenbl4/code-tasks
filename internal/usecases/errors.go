package usecases

import "errors"

var (
	// Ошибка, возвращаемая, если такой пользователь уже есть
	ErrUserAlreadyExists = errors.New("user already exists")
	// Ошибка, возвращаемая, если сессия уже есть
	ErrSessionAlreadyExists = errors.New("session already exists")
	// Ошибка, возвращаемая, если пароль неверен
	ErrIncorrectPassword = errors.New("incorrect password")
)
