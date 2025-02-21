package inmemstorage

import (
	"task-server/internal/domain"
	"task-server/internal/repository"

	"github.com/google/uuid"
)

// tasksStore - хранилище задач, конретная реализация интерфейса Storage, можеи быть заменена на другую реализацию
type tasksStore struct {
	// Хранилище задач в виде карты
	tasks map[uuid.UUID]domain.Task
}

// Создает новое хранилище задач
func NewTaskStore() *tasksStore { // либо Storage
	return &tasksStore{tasks: make(map[uuid.UUID]domain.Task)}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksStore) CreateTask() uuid.UUID {
	uuid := uuid.New()

	ts.tasks[uuid] = domain.Task{
		UUID:   uuid,
		Status: "in_progress",
	}

	return uuid
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksStore) GetTask(uuid uuid.UUID) (domain.Task, error) {
	t, ok := ts.tasks[uuid]
	if ok {
		return t, nil
	}
	return domain.Task{}, repository.ErrTaskNotFound
}

// Обновляет существующую задачу в хранилище
func (ts *tasksStore) UpdateTask(task domain.Task) error {
	_, ok := ts.tasks[task.UUID]
	if ok {
		ts.tasks[task.UUID] = task
		return nil
	}
	return repository.ErrTaskNotFound
}
