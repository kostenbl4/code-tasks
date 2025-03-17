package task

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"code-tasks/task-service/internal/usecases"
	"context"
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

	err := ts.sender.Send(ctx, task)
	if err == context.DeadlineExceeded {
		task.Result = domain.TaskResultError
		task.Stderr = "timeout"
		ts.repo.UpdateTask(task)
		return nil
	}

	return err
}

func (ts *tasksService) ListenTaskProcessor() error {

	tasks, err := ts.consumer.Consume()
	if err != nil {
		return err
	}

	// хз насколько рабочее решение в плане производительности, все таки один канал
	for task := range tasks {
		t, err := ts.repo.GetTask(task.UUID)
		if err != nil {
			return err
		}
		t.Status = "ready"
		t.Result = task.Result
		t.Stdout = task.Stdout
		t.Stderr = task.Stderr
		err = ts.repo.UpdateTask(t)
		if err != nil {
			return err
		}
	}

	return nil
}
