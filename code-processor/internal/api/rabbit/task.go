package rabbit

import (
	"code-processor/internal/domain"
	"code-processor/internal/usecases"
	"code-processor/pkg/broker"

	"encoding/json"
	"log"
)

var queueName = "code.process"

type RabbitHandler struct {
	client    broker.RabbitClient
	processor usecases.Processor
}

func NewRabbitHandler(consumeClient broker.RabbitClient, processor usecases.Processor) *RabbitHandler {

	_, err := consumeClient.CreateQueue(queueName, true, false)
	if err != nil {
		log.Fatal(err)
	}

	err = consumeClient.CreateBinding("code.process", "code.process", "code_requests")
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitHandler{
		client:    consumeClient,
		processor: processor,
	}
}

func (rh *RabbitHandler) ConsumeTasks() error {

	messages, err := rh.client.Consume(queueName, "code processor", false)
	if err != nil {
		return err
	}
	// канал для блокировки, временное решение
	blocking := make(chan struct{})
	go func() {
		for msg := range messages {
			log.Printf("new message: %v\n", msg.CorrelationId)

			var in domain.Task

			if err := json.Unmarshal(msg.Body, &in); err != nil {
				log.Printf("failed to unmarshall message: %v\n", err)
				continue
			}

			if err := rh.processor.Process(in); err != nil {
				log.Printf("failed to process message %v: %v\n", msg.MessageId, err)
				continue
			}

			log.Printf("payload: %v\n\n", in)
			if err := msg.Ack(false); err != nil {
				log.Println("Ack message failed ")
				continue
			}

			log.Printf("Acked and responded back to the msg %v\n", msg.MessageId)
		}
	}()
	log.Println("waiting for messages")
	<-blocking
	return nil
}
