package user

import (
	"task-server/internal/domain"
	"task-server/internal/repository"
	"task-server/internal/usecases"
	"task-server/utils"
)

// userService - сервис задач, конретная реализация интерфейса Task, может быть заменена на другую реализацию
type userService struct {
	repo      repository.User
	userCount int64
}

// Создает новый сервис пользователей
func NewUserService(repo repository.User) *userService { // либо Storage
	return &userService{repo: repo}
}

// Создает нового пользователя и добавляет его в хранилище
func (us *userService) RegisterUser(username, password string) (int64, error) {
	_, err := us.repo.GetUserByUsername(username)
	if err != repository.ErrUserNotFound {
		return -1, usecases.ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return -1, err
	}

	user := domain.User{
		ID:    us.userCount + 1,
		Login: username,
		Hpass: hashedPassword,
	}

	err = us.repo.CreateUser(user)
	if err != nil {
		return -1, err
	}
	us.userCount++

	return user.ID, nil
}

// Логин пользователя
func (us *userService) LoginUser(username, password string) (int64, error) {

	u, err := us.repo.GetUserByUsername(username)

	if err != nil {
		return -1, err
	}
	if !utils.CheckPassword(u.Hpass, password) {
		return -1, usecases.ErrIncorrectPassword
	}

	return u.ID, nil
}
