package main

import (
	"code-tasks/code-processor/internal/api/rabbit"
	"code-tasks/code-processor/internal/config"
	"code-tasks/code-processor/internal/usecases"
	httpsender "code-tasks/code-processor/internal/usecases/http_sender"
	"code-tasks/code-processor/internal/usecases/processor"
	pkgLogger "code-tasks/pkg/log"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	//pkgconfig "code-tasks/pkg/config"
	"net/http"

	//rabbitSender "code-processor/internal/usecases/rabbit_sender"
	"code-tasks/pkg/broker"
	"log"

	"github.com/docker/docker/client"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {

	var cfg config.Config
	

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Fatal(err)
	}

	logger, file := pkgLogger.NewLogger(cfg.Logger)
	slog.SetDefault(logger)
	
	defer file.Close()

	consumeConn, err := broker.ConnectRabbitMQ(cfg.Rabbit)
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

	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.45"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dockerClient.Close()

	codeExecutor := usecases.NewCodeExecutor(dockerClient)

	taskProcessor := processor.NewProcessor(resultSender, codeExecutor)
	rabbitHandler := rabbit.NewRabbitHandler(logger, consumeClient, taskProcessor)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := rabbitHandler.ConsumeTasks(ctx); err != nil {
		log.Fatal(err)
	}
}
