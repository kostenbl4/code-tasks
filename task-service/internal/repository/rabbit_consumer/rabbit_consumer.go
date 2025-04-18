package rabbitconsumer

import (
	"encoding/json"
	"log"

	"github.com/kostenbl4/code-tasks/pkg/broker"
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/repository"
)

var queueName = "code.results"

type rabbitmqConsumer struct {
	client broker.RabbitClient
}

func New(client broker.RabbitClient) repository.TaskConsumer {
	//client.CreateExchange("code_results", "direct", true, false)
	return rabbitmqConsumer{client: client}
}

func (rc rabbitmqConsumer) Consume() (<-chan domain.Task, error) {

	// queue, err := rc.client.CreateQueue("code.results", true, false)

	// if err != nil {
	// 	return nil, err
	// }

	// rc.client.CreateBinding(queue.Name, queue.Name, "code_results")

	messages, err := rc.client.Consume(queueName, "resulting", false)
	if err != nil {
		return nil, err
	}

	out := make(chan domain.Task, 100)

	var task domain.Task
	go func() {
		for msg := range messages {
			log.Printf("Got task %v result \n", msg.CorrelationId)

			if err := msg.Ack(false); err != nil {
				log.Println("Ack message failed")
				continue
			}
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Println(err)
				continue
			}
			out <- task
		}
	}()

	return out, nil
}
