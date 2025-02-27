package repository

import (
	"task-server/internal/domain"
)

// User - интерфейс для хранилища пользователей
type User interface {
	CreateUser(domain.User) error
	GetUserByUsername(string) (domain.User, error)
	//UpdateUser(domain.User) error
	//DeleteUser(int64) error
}
