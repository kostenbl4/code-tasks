package repository

import (
	"context"
	"task-service/internal/domain"
)

type TaskSender interface {
	Send(context.Context, domain.Task) error
}
