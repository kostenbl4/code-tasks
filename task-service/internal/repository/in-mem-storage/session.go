package inmemstorage

import (
	"context"
	"sync"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"
)

type sessionStore struct {
	sessions sync.Map
}

func NewSessionStore() repository.Session {
	return &sessionStore{}
}

func (ss *sessionStore) CreateSession(ctx context.Context, session domain.Session) error {
	ss.sessions.Store(session.SessionID, session)
	return nil
}

func (ss *sessionStore) GetSession(ctx context.Context, id string) (domain.Session, error) {
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

func (ss *sessionStore) GetSessionByUserId(ctx context.Context, userID int64) (domain.Session, bool) {
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

func (ss *sessionStore) DeleteSession(ctx context.Context, id string) error {
	ss.sessions.Delete(id)
	return nil
}
