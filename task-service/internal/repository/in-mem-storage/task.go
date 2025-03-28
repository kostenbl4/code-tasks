package inmemstorage

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"context"
	"sync"

	"github.com/google/uuid"
)

// tasksStore - хранилище задач в оперативной памяти
type tasksStore struct {
	// Хранилище задач в виде sync.Map
	tasks sync.Map
}

// Создает новое хранилище задач
func NewTaskStore() repository.Task {
	return &tasksStore{}
}

// Добавляет новую задачу в хранилище
func (ts *tasksStore) CreateTask(ctx context.Context, task domain.Task) error {
	ts.tasks.Store(task.ID, task)
	return nil
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksStore) GetTask(ctx context.Context, uuid uuid.UUID) (domain.Task, error) {
	if t, ok := ts.tasks.Load(uuid); ok {
		return t.(domain.Task), nil
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

// Обновляет существующую задачу в хранилище
func (ts *tasksStore) UpdateTask(ctx context.Context, task domain.Task) error {
	if _, ok := ts.tasks.Load(task.ID); ok {
		ts.tasks.Store(task.ID, task)
		return nil
	}
	return domain.ErrTaskNotFound
}
