package rabbit

import (
	"code-tasks/code-processor/internal/domain"
	"code-tasks/code-processor/internal/usecases"
	"code-tasks/pkg/broker"
	"context"
	"fmt"
	"log/slog"

	pkgLogger "code-tasks/pkg/log"

	"encoding/json"
)

var queueName = "code.process"

type RabbitHandler struct {
	logger *slog.Logger

	client    broker.RabbitClient
	processor usecases.Processor
}

func NewRabbitHandler(logger *slog.Logger, consumeClient broker.RabbitClient, processor usecases.Processor) *RabbitHandler {

	// _, err := consumeClient.CreateQueue(queueName, true, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = consumeClient.CreateBinding("code.process", "code.process", "code_requests")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &RabbitHandler{
		logger:    logger,
		client:    consumeClient,
		processor: processor,
	}
}

func (rh *RabbitHandler) ConsumeTasks(ctx context.Context) error {

	messages, err := rh.client.Consume(queueName, "code processor", false)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			log := rh.logger.With(
				slog.String("task_id", msg.CorrelationId),
				slog.String("routing_key", msg.RoutingKey),
			)

			var in domain.Task

			if err := json.Unmarshal(msg.Body, &in); err != nil {
				log.Error("failed to unmarshall message: ", pkgLogger.Error(err))
				continue
			}

			if err := rh.processor.Process(in); err != nil {
				log.Error("failed to process message: ", pkgLogger.Error(err))
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Error("Ack message failed: ", pkgLogger.Error(err))
				continue
			}

			log.Info("message completed")
		}
	}()
	rh.logger.Info("waiting for messages")

	<-ctx.Done()
	rh.logger.Info("shutting down the rabbit handler")
	err = rh.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close rabbit client: %w", err)
	}
	err = rh.processor.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop task processor: %w", err)
	}
	return nil
}
