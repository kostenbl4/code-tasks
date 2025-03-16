package processor

import (
	"code-processor/internal/domain"
	"code-processor/internal/usecases"
	"context"
	"time"
)

type processor struct {
	sender       usecases.Sender
	codeExecutor usecases.CodeExecutor
}

func NewProcessor(sender usecases.Sender, codeExecutor usecases.CodeExecutor) usecases.Processor {
	return processor{
		sender:       sender,
		codeExecutor: codeExecutor,
	}
}

func (p processor) Process(task domain.Task) error {

	// Задаем таймаут на выполнение кода
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	stdOut, stdErr, err := p.codeExecutor.Execute(ctx, task.Code, task.Translator)
	resTask := domain.Task{
		Translator: task.Translator,
		Code:       task.Code,
		UUID:       task.UUID,
		Status:     "ready",
	}
	if err != nil {
		resTask.Result = "error"
		resTask.Stderr = err.Error()
	}
	if stdErr != "" {
		resTask.Result = "error"
	} else {
		resTask.Result = "ok"
	}
	resTask.Stdout = stdOut
	resTask.Stderr = stdErr

	if err := p.sender.SendResult(ctx, resTask); err != nil {
		return err
	}

	return nil
}
