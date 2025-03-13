package repository

import "task-service/internal/domain"

type TaskConsumer interface{
	Consume() (<-chan domain.Task, error)
}