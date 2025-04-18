package repository

import (
	"context"

	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
)

type Session interface {
	CreateSession(context.Context, domain.Session) error
	GetSession(context.Context, string) (domain.Session, error)
	GetSessionByUserId(context.Context, int64) (domain.Session, bool)
	DeleteSession(context.Context, string) error
}
