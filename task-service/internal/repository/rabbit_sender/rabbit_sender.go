package rabbitsender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kostenbl4/code-tasks/pkg/broker"
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"

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
