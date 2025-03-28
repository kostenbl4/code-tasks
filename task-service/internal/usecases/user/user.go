package user

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"code-tasks/task-service/internal/usecases"
	"code-tasks/task-service/utils"
	"context"
	"log"
	"time"
)

// userService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type userService struct {
	repo repository.User

	defaultUserTimeout time.Duration
}

// Создает новый сервис пользователей
func NewUserService(repo repository.User) usecases.User {
	defaultUserTimeout := 5 * time.Second
	return &userService{
		repo: repo,

		defaultUserTimeout: defaultUserTimeout,
	}
}

// Создает нового пользователя и добавляет его в хранилище
func (us *userService) RegisterUser(username, password string) (int64, error) {

	log.Printf("Registering user %s", username)

	ctx, cancel := context.WithTimeout(context.Background(), us.defaultUserTimeout)
	defer cancel()

	log.Println("Getting user by username")
	_, err := us.repo.GetByUsername(ctx, username)

	if err != domain.ErrUserNotFound {
		return -1, domain.ErrUserAlreadyExists
	}

	log.Println("Hashing password")
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}

	user := domain.User{
		Username: username,
		Hpass:    hashedPassword,
	}

	log.Println("Creating user")
	id, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	log.Printf("User created with id %d", id)

	return id, nil
}

// Логин пользователя
func (us *userService) LoginUser(username, password string) (int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), us.defaultUserTimeout)
	defer cancel()
	u, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		return -1, err
	}
	if !utils.CheckPassword(u.Hpass, password) {
		return -1, domain.ErrIncorrectPassword
	}

	return u.ID, nil
}
