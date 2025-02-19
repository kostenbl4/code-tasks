package storage

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// Ошибка, возвращаемая, если задача не найдена
	errTaskNotFound = errors.New("task not found")
)

// Storage - интерфейс для хранилища задач
type Storage interface {
	// Создает новую задачу и возвращает её UUID
	CreateTask() uuid.UUID
	// Возвращает задачу по её UUID
	GetTask(uuid.UUID) (Task, error)
	// Обновляет существующую задачу
	UpdateTask(Task) error
}

// Task - структура задачи
type Task struct {
	UUID   uuid.UUID
	Status string
	Result []byte
}

// tasksStore - хранилище задач, конретная реализация интерфейса Storage, можеи быть заменена на другую реализацию
type tasksStore struct {
	// Хранилище задач в виде карты
	tasks map[uuid.UUID]Task
}

// Создает новое хранилище задач
func NewTaskStore() *tasksStore { // либо Storage
	return &tasksStore{tasks: make(map[uuid.UUID]Task)}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksStore) CreateTask() uuid.UUID {
	uuid := uuid.New()

	ts.tasks[uuid] = Task{
		UUID:   uuid,
		Status: "in_progress",
	}

	return uuid
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksStore) GetTask(uuid uuid.UUID) (Task, error) {
	t, ok := ts.tasks[uuid]
	if ok {
		return t, nil
	}
	return Task{}, errTaskNotFound
}

// Обновляет существующую задачу в хранилище
func (ts *tasksStore) UpdateTask(task Task) error {
	_, ok := ts.tasks[task.UUID]
	if ok {
		ts.tasks[task.UUID] = task
		return nil
	}
	return errTaskNotFound
}
