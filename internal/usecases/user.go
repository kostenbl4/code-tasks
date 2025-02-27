package usecases

import (
	//"task-server/internal/domain"
)

// User - интерфейс для сервиса пользователей
type User interface {
	RegisterUser(string, string) (int64, error)
	LoginUser(string, string) (int64, error)
	//UpdateUser(domain.User) error
	//DeleteUser(int64) error
}
