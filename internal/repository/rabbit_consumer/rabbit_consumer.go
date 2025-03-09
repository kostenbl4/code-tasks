package rabbitconsumer

import (
	"encoding/json"
	"log"

	"task-server/internal/domain"
	"task-server/internal/repository"
	"task-server/pkg/broker"
)

type rabbitmqConsumer struct {
	client broker.RabbitClient
}

func New(client broker.RabbitClient) repository.TaskConsumer {

	return rabbitmqConsumer{client: client}
}

func (rc rabbitmqConsumer) Consume() (<-chan domain.Task, error) {

	queue, err := rc.client.CreateQueue("code.process.results", true, false)

	if err != nil {
		return nil, err
	}

	rc.client.CreateBinding(queue.Name, queue.Name, "code_results")

	messages, err := rc.client.Consume(queue.Name, "resulting", false)
	if err != nil {
		return nil, err
	}

	out := make(chan domain.Task)

	var task domain.Task
	go func() {
		for msg := range messages {
			log.Printf("Got task %v result \n", msg.CorrelationId)

			if err := msg.Ack(false); err != nil {
				log.Println("Ack message failed")
				continue
			}

			log.Printf("result: %v", string(msg.Body))

			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Println(err)
				continue
			}
			out <- task
		}
	}()

	return out, nil
}
