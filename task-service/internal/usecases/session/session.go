package session

import (
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"code-tasks/task-service/internal/usecases"
	"code-tasks/task-service/utils"
	"context"
	"time"
)

const (
	SessionTokenLength = 32
)

type sessionManager struct {
	repo        repository.Session
	maxlifetime int64

	defaultSessionTimeout time.Duration
}

func NewSeessionManager(repo repository.Session, maxlifetime int64) usecases.Session {
	defaultSessionTimeout := 5 * time.Second
	return &sessionManager{
		repo:                  repo,
		maxlifetime:           maxlifetime,
		defaultSessionTimeout: defaultSessionTimeout,
	}
}

func (sm *sessionManager) CreateSession(userID int64) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), sm.defaultSessionTimeout)
	defer cancel()
	s, ok := sm.repo.GetSessionByUserId(ctx, userID)
	if ok {
		return s.SessionID, nil
	}

	sid, err := utils.GenerateSecureToken(SessionTokenLength)
	if err != nil {
		return "", err
	}

	s = domain.Session{UserID: userID, SessionID: sid}
	err = sm.repo.CreateSession(ctx, s)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (sm *sessionManager) GetSessionByID(sid string) (domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.defaultSessionTimeout)
	defer cancel()

	return sm.repo.GetSession(ctx, sid)
}

func (sm *sessionManager) DeleteSession(sid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), sm.defaultSessionTimeout)
	defer cancel()
	return sm.repo.DeleteSession(ctx, sid)
}
