package processor

import (
	"context"
	"time"

	"github.com/kostenbl4/code-tasks/code-processor/internal/domain"
	"github.com/kostenbl4/code-tasks/code-processor/internal/usecases"
)

type processor struct {
	sender       usecases.Sender
	codeExecutor usecases.CodeExecutor

	executeTimeout time.Duration
}

func NewProcessor(sender usecases.Sender, codeExecutor usecases.CodeExecutor) usecases.Processor {

	executeTimeout := 20 * time.Second
	return processor{
		sender:       sender,
		codeExecutor: codeExecutor,

		executeTimeout: executeTimeout,
	}
}

func (p processor) Process(task domain.Task) error {

	// Задаем таймаут на выполнение кода
	ctx, cancel := context.WithTimeout(context.Background(), p.executeTimeout)
	defer cancel()
	stdOut, stdErr, err := p.codeExecutor.Execute(ctx, task.Code, task.Translator)
	task.Status = "ready"
	if err != nil {
		task.Result = "error"
		task.Stderr = err.Error()
	}
	if stdErr != "" {
		task.Result = "error"
	} else {
		task.Result = "ok"
	}
	task.Stdout = stdOut
	task.Stderr = stdErr

	if err := p.sender.SendResult(ctx, task); err != nil {
		return err
	}

	return nil
}

func (p processor) Stop() error {
	p.executeTimeout = 30 * time.Second
	return nil
}
