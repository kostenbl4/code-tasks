package main

import (
	"code-processor/pkg/broker"
	"code-processor/processor"
	"code-processor/types"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/client"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	consumeConn, err := broker.ConnectRabbitMQ("myuser", "mypassword", "rabbitmq:5672", "")
	if err != nil {
		log.Fatal(err)
	}
	defer consumeConn.Close()

	consumeClient, err := broker.NewRabbitClient(consumeConn)
	if err != nil {
		log.Fatal(err)
	}
	defer consumeClient.Close()

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

	sendConn, err := broker.ConnectRabbitMQ("myuser", "mypassword", "rabbitmq:5672", "")
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

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.45"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = processor.LoadImages(cli)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range messages {
			log.Printf("new message: %v\n", msg.CorrelationId)

			var in types.ProcessTask
			var out types.ProcessTask
			if err := json.Unmarshal(msg.Body, &in); err != nil {
				log.Printf("failed to unmarshall message: %v\n", err)
				continue
			}

			out = types.ProcessTask{
				Translator: in.Translator,
				Code:       in.Code,
				UUID:       in.UUID,
				Status:     "ready",
			}

			log.Printf("payload: %v\n\n", in)
			if err := msg.Ack(false); err != nil {
				log.Println("Ack message failed ")
				continue
			}

			var stdOut string
			var stdErr string

			stdOut, stdErr, err := processor.ExecuteCode(cli, in.Code, in.Translator)

			if err != nil {
				log.Printf("Error executing code in message %s: %s\n", msg.MessageId, err)
				out.Result = "error"
				out.Stderr = err.Error()
			}
			if stdErr != "" {
				out.Result = "error"
			} else {
				out.Result = "ok"
			}
			out.Stdout = stdOut
			out.Stderr = stdErr

			log.Printf("task responce: %v\n", out)

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
