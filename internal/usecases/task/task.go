package service

import (
	"task-server/internal/domain"
	"task-server/internal/repository"

	"github.com/google/uuid"
)

// tasksService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type tasksService struct {
	repo repository.Task
}

// Создает новое хранилище задач
func NewTaskService(repo repository.Task) *tasksService { // либо Storage
	return &tasksService{repo: repo}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksService) CreateTask() uuid.UUID {
	return ts.repo.CreateTask()
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksService) GetTask(uuid uuid.UUID) (domain.Task, error) {
	return ts.repo.GetTask(uuid)
}

// Обновляет существующую задачу в хранилище
func (ts *tasksService) UpdateTask(task domain.Task) error {
	return ts.repo.UpdateTask(task)
}
