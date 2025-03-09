package repository

import "task-server/internal/domain"

type TaskConsumer interface{
	Consume() (<-chan domain.Task, error)
}