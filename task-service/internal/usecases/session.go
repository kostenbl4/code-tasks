package usecases

import (
	"code-tasks/task-service/internal/domain"
)

// Session - интерфейс для сервиса сессий
type Session interface {
	// Создает новую сессию для пользователя
	CreateSession(int64) (string, error)
	// Возвращает сессию по её ID
	GetSessionByID(string) (domain.Session, error)
	// Удаляет сессию по её ID
	DeleteSession(string) error
	// Проверяет, действительна ли сессия
}
