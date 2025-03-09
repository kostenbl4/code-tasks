package task

import (
	"context"
	"task-server/internal/domain"
	"task-server/internal/repository"
	"task-server/internal/usecases"
	"time"

	"github.com/google/uuid"
)

// tasksService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type tasksService struct {
	repo     repository.Task
	sender   repository.TaskSender
	consumer repository.TaskConsumer
}

// Создает новый сервис задач
func NewTaskService(repo repository.Task, sender repository.TaskSender, consumer repository.TaskConsumer) usecases.Task {
	return &tasksService{repo: repo, sender: sender, consumer: consumer}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksService) CreateTask(translator, code string) (domain.Task, error) {
	id := uuid.New()
	// проверка на то, что нового uuid нет среди существующих
	for {
		if _, err := ts.repo.GetTask(id); err == domain.ErrTaskNotFound {
			break
		}
		id = uuid.New()
	}

	task := domain.Task{
		Translator: translator,
		Code:       code,
		UUID:       id,
		Status:     "in_progress",
	}

	ts.repo.CreateTask(task)

	return task, nil
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksService) GetTask(uuid uuid.UUID) (domain.Task, error) {
	return ts.repo.GetTask(uuid)
}

// Обновляет существующую задачу в хранилище
func (ts *tasksService) UpdateTask(task domain.Task) error {
	return ts.repo.UpdateTask(task)
}

func (ts *tasksService) SendTask(task domain.Task) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ts.sender.Send(ctx, task)
}

func (ts *tasksService) ListenTaskProcessor() error {

	tasks, err := ts.consumer.Consume()
	if err != nil {
		return err
	}

	for task := range tasks {
		t, err := ts.repo.GetTask(task.UUID)
		if err != nil {
			return err
		}
		t.Status = "ready"
		t.Stdout = task.Stdout
		t.Stderr = task.Stderr
		err = ts.repo.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}
