package usecases

import (
	"context"

	"github.com/kostenbl4/code-tasks/code-processor/internal/domain"
)

type Sender interface {
	SendResult(ctx context.Context, task domain.Task) error
}
