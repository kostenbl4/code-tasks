package repository

import (
	"code-tasks/task-service/internal/domain"
	"context"
)

// User - интерфейс для хранилища пользователей
type User interface {
	CreateUser(context.Context, domain.User) (int64, error)
	GetByUsername(context.Context, string) (domain.User, error)
	//UpdateUser(domain.User) error
	//DeleteUser(int64) error
}
