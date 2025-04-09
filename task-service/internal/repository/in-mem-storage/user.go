package inmemstorage

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"
)

// userStore - хранилище пользователей в оперативной памяти
type userStore struct {
	// Хранилище пользователей в виде sync.Map, на данном этапе будем хранить ключи в виде username
	users     sync.Map
	userCount atomic.Int64
}

// Создает новое хранилище пользователей
func NewUserStore() repository.User {
	return &userStore{}
}

// Создает новую задачу и добавляет её в хранилище
func (us *userStore) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	user.ID = us.userCount.Load()
	us.users.Store(user.Username, user)
	us.userCount.Add(1)
	return user.ID, nil
}

func (us *userStore) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	if u, ok := us.users.Load(username); ok {
		user, ok := u.(domain.User)
		if !ok {
			return domain.User{}, domain.ErrUserNotFound
		}
		return user, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}
