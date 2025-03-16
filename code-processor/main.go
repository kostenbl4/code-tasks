package main

import (
	"code-processor/internal/api/rabbit"
	"code-processor/internal/usecases"
	httpsender "code-processor/internal/usecases/http_sender"
	"code-processor/internal/usecases/processor"
	"net/http"

	//rabbitSender "code-processor/internal/usecases/rabbit_sender"
	"code-processor/pkg/broker"
	"log"

	"github.com/docker/docker/client"
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

	// Есть 2 варианта отправки: rabbitmq или http

	// sendConn, err := broker.ConnectRabbitMQ("myuser", "mypassword", "rabbitmq:5672", "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// sendClient, err := broker.NewRabbitClient(sendConn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// resultSender := rabbitSender.NewRabbitSender(sendClient)

	httpClient := http.Client{}
	resultSender := httpsender.NewHttpSender(httpClient)

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.45"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	codeExecutor := usecases.NewCodeExecutor(cli)

	taskProcessor := processor.NewProcessor(resultSender, codeExecutor)
	rabbitHandler := rabbit.NewRabbitHandler(consumeClient, taskProcessor)

	if err := rabbitHandler.ConsumeTasks(); err != nil {
		log.Fatal(err)
	}

	// TODO: graceful shutdown
}
