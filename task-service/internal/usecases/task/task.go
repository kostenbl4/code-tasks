package task

import (
	"context"
	"time"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"
	"github.com/kostenbl4/code-tasks/task-service/internal/usecases"

	"github.com/google/uuid"
)

// tasksService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type tasksService struct {
	repo     repository.Task
	sender   repository.TaskSender
	consumer repository.TaskConsumer

	defaultTaskTimeout time.Duration
	sendTaskTimeout    time.Duration
}

// Создает новый сервис задач
func NewTaskService(repo repository.Task, sender repository.TaskSender, consumer repository.TaskConsumer) usecases.Task {
	defaultTaskTimeout := 5 * time.Second
	sendTaskTimeout := 5 * time.Second
	return &tasksService{
		repo:               repo,
		sender:             sender,
		consumer:           consumer,
		defaultTaskTimeout: defaultTaskTimeout,
		sendTaskTimeout:    sendTaskTimeout,
	}
}

// Создает новую задачу и добавляет её в хранилище
func (ts *tasksService) CreateTask(translator, code string, userID int64) (domain.Task, error) {
	id := uuid.New()
	// проверка на то, что нового uuid нет среди существующих
	for {
		ctx, _ := context.WithTimeout(context.Background(), ts.defaultTaskTimeout)
		if _, err := ts.repo.GetTask(ctx, id); err == domain.ErrTaskNotFound {
			break
		}
		id = uuid.New()
	}

	task := domain.Task{
		Translator: translator,
		Code:       code,
		ID:         id,
		UserID:     userID,
		Status:     "in_progress",
	}

	ctx, cancel := context.WithTimeout(context.Background(), ts.defaultTaskTimeout)
	defer cancel()
	ts.repo.CreateTask(ctx, task)

	return task, nil
}

// Возвращает задачу по её UUID, если она существует
func (ts *tasksService) GetTask(uuid uuid.UUID) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ts.defaultTaskTimeout)
	defer cancel()
	return ts.repo.GetTask(ctx, uuid)
}

// Обновляет существующую задачу в хранилище
func (ts *tasksService) UpdateTask(task domain.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), ts.defaultTaskTimeout)
	defer cancel()
	return ts.repo.UpdateTask(ctx, task)
}

func (ts *tasksService) SendTask(task domain.Task) error {

	ctx, cancel := context.WithTimeout(context.Background(), ts.sendTaskTimeout)
	defer cancel()

	err := ts.sender.Send(ctx, task)
	if err == context.DeadlineExceeded {
		task.Result = domain.TaskResultError
		task.Stderr = "timeout"
		ts.repo.UpdateTask(ctx, task)
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

		ctx, cancel := context.WithTimeout(context.Background(), ts.defaultTaskTimeout)
		defer cancel()

		t, err := ts.repo.GetTask(ctx, task.ID)
		if err != nil {
			return err
		}
		t.Status = "ready"
		t.Result = task.Result
		t.Stdout = task.Stdout
		t.Stderr = task.Stderr
		err = ts.repo.UpdateTask(ctx, t)
		if err != nil {
			return err
		}
	}

	return nil
}
