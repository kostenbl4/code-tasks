package usecases

import (
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"

	"github.com/google/uuid"
)

// Task - интерфейс для сервиса задач
type Task interface {
	// Создает новую задачу и возвращает её UUID
	CreateTask(string, string, int64) (domain.Task, error)
	// Возвращает задачу по её UUID
	GetTask(uuid.UUID) (domain.Task, error)
	// Обновляет существующую задачу
	UpdateTask(domain.Task) error

	SendTask(domain.Task) error

	ListenTaskProcessor() error
}
