package repository

import "github.com/kostenbl4/code-tasks/task-service/internal/domain"

type TaskConsumer interface {
	Consume() (<-chan domain.Task, error)
}
