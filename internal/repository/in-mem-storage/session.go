package inmemstorage

import (
	"sync"
	"task-server/internal/domain"
	"task-server/internal/repository"
)

type sessionStore struct {
	sessions sync.Map
}

func NewSessionStore() repository.Session {
	return &sessionStore{}
}

func (ss *sessionStore) CreateSession(session domain.Session) error {
	ss.sessions.Store(session.SessionID, session)
	return nil
}

func (ss *sessionStore) GetSession(id string) (domain.Session, error) {
	value, ok := ss.sessions.Load(id)
	if !ok {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	session, ok := value.(domain.Session)
	if !ok {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	return session, nil
}

func (ss *sessionStore) GetSessionByUserId(userID int64) (domain.Session, bool) {
	var s domain.Session
	found := false
	ss.sessions.Range(func(key, value interface{}) bool {
		session, ok := value.(domain.Session)
		if !ok {
			return false
		}
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
