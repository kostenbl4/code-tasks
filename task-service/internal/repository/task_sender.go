package repository

import (
	"context"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
)

type TaskSender interface {
	Send(context.Context, domain.Task) error
}
