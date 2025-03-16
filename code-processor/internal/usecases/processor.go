package usecases

import (
	"code-processor/internal/domain"
	"context"
	"time"
)

type Processor interface {
	Process(task domain.Task) error
}

type processor struct {
	sender       Sender
	codeExecutor CodeExecutor
}

func NewProcessor(sender Sender, codeExecutor CodeExecutor) Processor {
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
