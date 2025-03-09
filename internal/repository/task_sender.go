package repository

import (
	"context"
	"task-server/internal/domain"
)

type TaskSender interface {
	Send(context.Context, domain.Task) error
}
