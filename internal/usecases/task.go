package usecases

import (
	"task-server/internal/domain"

	"github.com/google/uuid"
)

// Task - интерфейс для сервиса задач
type Task interface {
	// Создает новую задачу и возвращает её UUID
	CreateTask(string, string) (domain.Task, error)
	// Возвращает задачу по её UUID
	GetTask(uuid.UUID) (domain.Task, error)
	// Обновляет существующую задачу
	UpdateTask(domain.Task) error

	SendTask(domain.Task) error

	ListenTaskProcessor() error
}
