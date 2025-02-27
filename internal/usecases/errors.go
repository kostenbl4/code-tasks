package usecases

import "errors"

var (
	// Ошибка, возвращаемая, если такой пользователь уже есть
	ErrUserAlreadyExists = errors.New("user already exists")
	// Ошибка, возвращаемая, если пароль неверен
	ErrIncorrectPassword = errors.New("incorrect password")
)
