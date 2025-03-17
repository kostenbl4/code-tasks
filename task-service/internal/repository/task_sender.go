package repository

import (
	"context"
	"code-tasks/task-service/internal/domain"
)

type TaskSender interface {
	Send(context.Context, domain.Task) error
}
