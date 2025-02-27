package inmemstorage

import (
	"sync"
	"task-server/internal/domain"
	"task-server/internal/repository"
)

// userStore - хранилище пользователей в оперативной памяти
type userStore struct {
	// Хранилище пользователей в виде sync.Map, на данном этапе будем хранить ключи в виде username
	users sync.Map
}

// Создает новое хранилище пользователей
func NewUserStore() *userStore {
	return &userStore{}
}

// Создает новую задачу и добавляет её в хранилище
func (us *userStore) CreateUser(user domain.User) error {
	us.users.Store(user.Login, user)
	return nil
}

func (us *userStore) GetUserByUsername(username string) (domain.User, error) {
	if u, ok := us.users.Load(username); ok {
		return u.(domain.User), nil
	}
	return domain.User{}, repository.ErrUserNotFound
}
