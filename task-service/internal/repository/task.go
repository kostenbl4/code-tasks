package repository

import (
	"code-tasks/task-service/internal/domain"

	"github.com/google/uuid"
)

// Task - интерфейс для хранилища задач
type Task interface {
	CreateTask(domain.Task)
	GetTask(uuid.UUID) (domain.Task, error)
	UpdateTask(domain.Task) error
}



