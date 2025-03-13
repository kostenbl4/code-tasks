package domain

import "errors"

var (
	ErrBadRequest = errors.New("bad request")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrIncorrectPassword = errors.New("incorrect password")

	ErrTaskNotFound    = errors.New("task not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)
