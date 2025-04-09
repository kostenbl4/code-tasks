package repository

import (
	"context"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"

	"github.com/google/uuid"
)

// Task - интерфейс для хранилища задач
type Task interface {
	CreateTask(context.Context, domain.Task) error
	GetTask(context.Context, uuid.UUID) (domain.Task, error)
	UpdateTask(context.Context, domain.Task) error
}
