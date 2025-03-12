package rabbitsender

import (
	"context"
	"encoding/json"
	"log"
	"task-server/internal/domain"
	"task-server/internal/repository"
	rabbitconsumer "task-server/internal/repository/rabbit_consumer"
	"task-server/pkg/broker"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqSender struct {
	client broker.RabbitClient
}

func New(client broker.RabbitClient) repository.TaskSender {
	client.CreateExchange("code_requests", "direct", true, false)
	return rabbitmqSender{client: client}
}

func (rs rabbitmqSender) Send(ctx context.Context, task domain.Task) error {
	log.Printf("sending message: taskID = %v", task.UUID)
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return rs.client.Send(ctx, "code_requests", "code.process", amqp.Publishing{
		ContentType:   "application/json",
		DeliveryMode:  amqp.Persistent,
		ReplyTo:       rabbitconsumer.ReplyQueue,
		CorrelationId: task.UUID.String(),
		Body:          data,
	})
}
