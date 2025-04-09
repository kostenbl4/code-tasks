package repository

import (
	"context"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
)

// User - интерфейс для хранилища пользователей
type User interface {
	CreateUser(context.Context, domain.User) (int64, error)
	GetByUsername(context.Context, string) (domain.User, error)
	//UpdateUser(domain.User) error
	//DeleteUser(int64) error
}
