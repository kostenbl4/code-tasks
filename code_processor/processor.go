package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"task-server/code_processor/types"
	"task-server/pkg/broker"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	consumeConn, err := broker.ConnectRabbitMQ("guest", "guest", "localhost:5672", "")
	if err != nil {
		log.Fatal(err)
	}
	defer consumeConn.Close()

	consumeClient, err := broker.NewRabbitClient(consumeConn)
	if err != nil {
		log.Fatal(err)
	}
	defer consumeClient.Close()

	consumeClient.CreateExchange("code_results", "direct", true, false)

	queue, err := consumeClient.CreateQueue("code.process", true, false)
	if err != nil {
		log.Fatal(err)
	}

	err = consumeClient.CreateBinding("code.process", "code.process", "code_requests")
	if err != nil {
		log.Fatal(err)
	}

	messages, err := consumeClient.Consume(queue.Name, "code processor", false)
	if err != nil {
		log.Fatal(err)
	}

	sendConn, err := broker.ConnectRabbitMQ("guest", "guest", "localhost:5672", "")
	if err != nil {
		log.Fatal(err)
	}
	sendClient, err := broker.NewRabbitClient(sendConn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var blocking chan struct{}

	go func() {
		for msg := range messages {
			log.Printf("new message: %v\n", msg)

			var in types.ProcessTask
			if err := json.Unmarshal(msg.Body, &in); err != nil {
				log.Printf("failed to unmarshall message: %v\n", err)
				msg.Nack(false, true)
				continue
			}

			log.Printf("payload: %v\n\n", in)
			if err := msg.Ack(false); err != nil {
				log.Println("Ack message failed ")
				continue
			}

			time.Sleep(2 * time.Second)

			out := types.ProcessTask{
				Translator: in.Translator,
				Code:       in.Code,
				UUID:       in.UUID,
				Status:     "ready",
				Stdout:     "some stdout",
				Stderr:     "some stderr",
			}
			data, err := json.Marshal(out)
			if err != nil {
				log.Fatal(err)
			}
			if err := sendClient.Send(ctx, "code_results", msg.ReplyTo, amqp.Publishing{
				ContentType:   "application/json",
				DeliveryMode:  amqp.Persistent,
				Body:          data,
				CorrelationId: msg.CorrelationId,
			}); err != nil {
				log.Fatal(err)
			}

			log.Printf("Acked and responded back to the msg %v\n", msg.MessageId)
		}
	}()
	fmt.Println("waiting for messages")
	<-blocking
}
