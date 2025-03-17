package repository

import "code-tasks/task-service/internal/domain"

type TaskConsumer interface{
	Consume() (<-chan domain.Task, error)
}