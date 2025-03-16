package usecases

import (
	"code-processor/internal/domain"
)

type Processor interface {
	Process(task domain.Task) error
}

