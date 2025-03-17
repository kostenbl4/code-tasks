package inmemstorage

import (
	"sync"
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"

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
func (ts *tasksStore) CreateTask(task domain.Task) {
	ts.tasks.Store(task.UUID, task)
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksStore) GetTask(uuid uuid.UUID) (domain.Task, error) {
	if t, ok := ts.tasks.Load(uuid); ok {
		return t.(domain.Task), nil
	}
	return domain.Task{}, domain.ErrTaskNotFound
}

// Обновляет существующую задачу в хранилище
func (ts *tasksStore) UpdateTask(task domain.Task) error {
	if _, ok := ts.tasks.Load(task.UUID); ok {
		ts.tasks.Store(task.UUID, task)
		return nil
	}
	return domain.ErrTaskNotFound
}
