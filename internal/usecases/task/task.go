package task

import (
	"task-server/internal/domain"
	"task-server/internal/repository"
	"task-server/internal/usecases"

	"github.com/google/uuid"
)

// tasksService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type tasksService struct {
	repo repository.Task
}

// Создает новый сервис задач
func NewTaskService(repo repository.Task)  usecases.Task{ 
	return &tasksService{repo: repo}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksService) CreateTask() uuid.UUID {
	id := uuid.New()
	// проверка на то, что нового uuid нет среди существующих
	for {
		if _, err := ts.repo.GetTask(id); err == domain.ErrTaskNotFound {
			break
		}
		id = uuid.New()
	}

	ts.repo.CreateTask(domain.Task{
		UUID:   id,
		Status: "in_progress",
	})
	return id
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksService) GetTask(uuid uuid.UUID) (domain.Task, error) {
	return ts.repo.GetTask(uuid)
}

// Обновляет существующую задачу в хранилище
func (ts *tasksService) UpdateTask(task domain.Task) error {
	return ts.repo.UpdateTask(task)
}
