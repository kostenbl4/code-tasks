package inmemstorage

import (
	"sync"
	"task-server/internal/domain"
	"task-server/internal/repository"

	"github.com/google/uuid"
)

// tasksStore - хранилище задач, конретная реализация интерфейса Storage, можеи быть заменена на другую реализацию
type tasksStore struct {
	// Хранилище задач в виде sync.Map
	tasks sync.Map
}

// Создает новое хранилище задач
func NewTaskStore() *tasksStore { // либо Storage
	return &tasksStore{}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksStore) CreateTask() uuid.UUID {

	id := uuid.New()
	// проверка на то, что нового uuid нет среди существующих
	for {
		if _, ok := ts.tasks.Load(id); !ok {
			break
		}
		id = uuid.New()
	}

	ts.tasks.Store(id, domain.Task{
		UUID:   id,
		Status: "in_progress",
	})

	return id
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksStore) GetTask(uuid uuid.UUID) (domain.Task, error) {
	if t, ok := ts.tasks.Load(uuid); ok {
		return t.(domain.Task), nil
	}
	return domain.Task{}, repository.ErrTaskNotFound
}

// Обновляет существующую задачу в хранилище
func (ts *tasksStore) UpdateTask(task domain.Task) error {
	if _, ok := ts.tasks.Load(task.UUID); ok {
		ts.tasks.Store(task.UUID, task)
		return nil
	}
	return repository.ErrTaskNotFound
}
