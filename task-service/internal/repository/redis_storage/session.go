package redisstorage

import (
	"context"
	"time"

	"github.com/kostenbl4/code-tasks/pkg/cache"
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"
)

type sessionStore struct {
	cache cache.Cache
	ttl   time.Duration
}

func NewSessionStore(cache cache.Cache, ttl time.Duration) repository.Session {
	return &sessionStore{
		cache: cache,
		ttl:   ttl,
	}
}

func (s *sessionStore) CreateSession(ctx context.Context, session domain.Session) error {
	return s.cache.Set(ctx, session.SessionID, session, s.ttl)
}

func (s *sessionStore) GetSession(ctx context.Context, id string) (domain.Session, error) {
	var session domain.Session
	err := s.cache.Get(ctx, id, &session)
	return session, err
}

func (s *sessionStore) GetSessionByUserId(ctx context.Context, userID int64) (domain.Session, bool) {
	return domain.Session{}, false
}

func (s *sessionStore) DeleteSession(ctx context.Context, id string) error {
	return s.cache.Delete(ctx, id)
}
