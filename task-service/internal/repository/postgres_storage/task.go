package postgresstorage

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskStore struct {
	pool *pgxpool.Pool
}

func NewTaskStore(pool *pgxpool.Pool) repository.Task {
	return &taskStore{pool: pool}
}

func (ts *taskStore) CreateTask(ctx context.Context, task domain.Task) error {
	_, err := ts.pool.Exec(ctx, `INSERT INTO tasks (id, user_id, translator, code, task_status, result, stdout, stderr) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		task.ID,
		task.UserID,
		task.Translator,
		task.Code,
		task.Status,
		task.Result,
		task.Stdout,
		task.Stderr,
	)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	return nil
}

func (ts *taskStore) GetTask(ctx context.Context, uuid uuid.UUID) (domain.Task, error) {
	var task domain.Task
	err := ts.pool.QueryRow(ctx, `SELECT id, translator, code, task_status, result, stdout, stderr FROM tasks WHERE id = $1`, uuid).Scan(
		&task.ID,
		&task.Translator,
		&task.Code,
		&task.Status,
		&task.Result,
		&task.Stdout,
		&task.Stderr,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (ts *taskStore) UpdateTask(ctx context.Context, task domain.Task) error {
	_, err := ts.pool.Exec(ctx, `UPDATE tasks SET task_status = $1, result = $2, stdout = $3, stderr = $4 WHERE id = $5`,
		task.Status,
		task.Result,
		task.Stdout,
		task.Stderr,
		task.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrTaskNotFound
		}
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil

}
