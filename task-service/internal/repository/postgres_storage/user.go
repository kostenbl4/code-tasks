package postgresstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userStore struct {
	pool *pgxpool.Pool
}

func NewUserStore(pool *pgxpool.Pool) repository.User {
	return &userStore{pool: pool}
}

func (us *userStore) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	var id int64
	err := us.pool.QueryRow(ctx, `INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id`, user.Username, user.Hpass).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, domain.ErrBadRequest
		}
		return -1, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (us *userStore) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	err := us.pool.QueryRow(ctx, `SELECT id, username, hashed_password FROM users WHERE username = $1`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Hpass,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
