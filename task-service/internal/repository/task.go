package repository

import (
	"code-tasks/task-service/internal/domain"
	"context"

	"github.com/google/uuid"
)

// Task - интерфейс для хранилища задач
type Task interface {
	CreateTask(context.Context, domain.Task) error
	GetTask(context.Context, uuid.UUID) (domain.Task, error)
	UpdateTask(context.Context, domain.Task) error
}
