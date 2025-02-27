package session

import (
	"task-server/internal/domain"
	"task-server/internal/repository"
	"task-server/utils"
)

type sessionManager struct {
	repo        repository.Session
	maxlifetime int64
}

func NewSeessionManager(repo repository.Session, maxlifetime int64) *sessionManager {
	return &sessionManager{repo: repo, maxlifetime: maxlifetime}
}

func (sm *sessionManager) CreateSession(userID int64) (string, error) {

	s, ok := sm.repo.GetSessionByUserId(userID)
	if ok {
		return s.SessionID, nil
	}

	sid, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}

	s = domain.Session{UserID: userID, SessionID: sid}
	err = sm.repo.CreateSession(s)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (sm *sessionManager) GetSessionByID(sid string) (domain.Session, error) {
	return sm.repo.GetSession(sid)
}

func (sm *sessionManager) DeleteSession(sid string) error {
	return sm.repo.DeleteSession(sid)
}
