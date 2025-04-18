package rabbitsender

import (
	"context"
	"encoding/json"

	"github.com/kostenbl4/code-tasks/code-processor/internal/domain"
	"github.com/kostenbl4/code-tasks/pkg/broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitSender - вариант отправки результата в rabbit
type RabbitSender struct {
	client broker.RabbitClient
}

func NewRabbitSender(client broker.RabbitClient) *RabbitSender {
	return &RabbitSender{client: client}
}

func (rs RabbitSender) SendResult(ctx context.Context, task domain.Task) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	if err := rs.client.Send(ctx, "code_results", "code.results", amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         data,
	}); err != nil {
		return err
	}
	return nil
}
