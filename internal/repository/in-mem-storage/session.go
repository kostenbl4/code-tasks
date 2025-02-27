package inmemstorage

import (
	"sync"
	"task-server/internal/domain"
	"task-server/internal/repository"
)

type sessionStore struct {
	sessions sync.Map
}

func NewSessionStore() *sessionStore {
	return &sessionStore{}
}

func (ss *sessionStore) CreateSession(session domain.Session) error {
	ss.sessions.Store(session.SessionID, session)
	return nil
}

func (ss *sessionStore) GetSession(id string) (domain.Session, error) {
	value, ok := ss.sessions.Load(id)
	if !ok {
		return domain.Session{}, repository.ErrSessionNotFound
	}
	return value.(domain.Session), nil
}

func (ss *sessionStore) GetSessionByUserId(userID int64) (domain.Session, bool) {
	var s domain.Session
	found := false
	ss.sessions.Range(func(key, value interface{}) bool {
		session := value.(domain.Session)
		if session.UserID == userID {
			s = session
			found = true
		}
		return true
	})
	return s, found
}

func (ss *sessionStore) DeleteSession(id string) error {
	ss.sessions.Delete(id)
	return nil
}


