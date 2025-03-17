package usecases

import (
	"code-tasks/code-processor/internal/domain"
	"context"

)


type Sender interface {
	SendResult(ctx context.Context, task domain.Task) error
}