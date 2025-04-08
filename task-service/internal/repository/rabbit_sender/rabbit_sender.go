package rabbitsender

import (
	"code-tasks/pkg/broker"
	"code-tasks/task-service/internal/domain"
	"code-tasks/task-service/internal/repository"
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqSender struct {
	client broker.RabbitClient
}

func New(client broker.RabbitClient) repository.TaskSender {
	//client.CreateExchange("code_requests", "direct", true, false)
	return rabbitmqSender{client: client}
}

func (rs rabbitmqSender) Send(ctx context.Context, task domain.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to send message with taskID = %v : %w", task.ID, err)
	}
	return rs.client.Send(ctx, "code_requests", "code.process", amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		CorrelationId: task.ID.String(),
		Body:          data,
	})
}
