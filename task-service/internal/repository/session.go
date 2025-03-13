package repository

import "task-service/internal/domain"

type Session interface {
	CreateSession(domain.Session) error
	GetSession(string) (domain.Session, error)
	GetSessionByUserId(int64) (domain.Session, bool)
	DeleteSession(string) error
}
