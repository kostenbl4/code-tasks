package usecases

import (
	"code-tasks/code-processor/internal/domain"
)

type Processor interface {
	Process(task domain.Task) error
	Stop() error
}

