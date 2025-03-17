package user

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"code-tasks/task-service/internal/usecases"
	"code-tasks/task-service/utils"
)

// userService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type userService struct {
	repo repository.User
}

// Создает новый сервис пользователей
func NewUserService(repo repository.User) usecases.User {
	return &userService{repo: repo}
}

// Создает нового пользователя и добавляет его в хранилище
func (us *userService) RegisterUser(username, password string) (int64, error) {
	_, err := us.repo.GetUserByUsername(username)
	if err != domain.ErrUserNotFound {
		return -1, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}

	user := domain.User{
		Login: username,
		Hpass: hashedPassword,
	}

	err = us.repo.CreateUser(user)
	if err != nil {
		return -1, err
	}

	return user.ID, nil
}

// Логин пользователя
func (us *userService) LoginUser(username, password string) (int64, error) {

	u, err := us.repo.GetUserByUsername(username)
	if err != nil {
		return -1, err
	}
	if !utils.CheckPassword(u.Hpass, password) {
		return -1, domain.ErrIncorrectPassword
	}

	return u.ID, nil
}
