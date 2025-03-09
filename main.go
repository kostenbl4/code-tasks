package main

import (
	"flag"
	"log"
	"task-server/internal/api/http"
	inmemstorage "task-server/internal/repository/in-mem-storage"
	rabbitconsumer "task-server/internal/repository/rabbit_consumer"
	rabbitsender "task-server/internal/repository/rabbit_sender"
	"task-server/internal/usecases/session"
	"task-server/internal/usecases/task"
	"task-server/internal/usecases/user"
	"task-server/pkg/broker"
	httpServer "task-server/pkg/http"

	_ "task-server/docs"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title Task-server API
// @version 1.0
// @description This is a sample server.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {

	addr := flag.String("addr", ":8080", "address for the http server")

	// Создаем конфигурацию для API сервера
	config := httpServer.Config{
		Addr: *addr,
	}

	// Создаем хранилище, менеджер, сессий
	sessionStore := inmemstorage.NewSessionStore()
	sessionManager := session.NewSeessionManager(sessionStore, 3600)

	// Создаем хранилище, сервис, обработчик задач
	taskStore := inmemstorage.NewTaskStore()

	sendConn, err := broker.ConnectRabbitMQ("guest", "guest", "localhost:5672", "")
	if err != nil {
		log.Fatal(err)
	}
	sendClient, err := broker.NewRabbitClient(sendConn)
	if err != nil {
		log.Fatal(err)
	}
	taskSender := rabbitsender.New(sendClient)

	consumeConn, err := broker.ConnectRabbitMQ("guest", "guest", "localhost:5672", "")
	if err != nil {
		log.Fatal(err)
	}
	consumeClient, err := broker.NewRabbitClient(consumeConn)
	if err != nil {
		log.Fatal(err)
	}
	taskConsumer := rabbitconsumer.New(consumeClient)

	taskService := task.NewTaskService(taskStore, taskSender, taskConsumer)
	taskHandler := http.NewTaskHandler(taskService, sessionManager)

	go taskService.ListenTaskProcessor()

	// Создаем хранилище, сервис, обработчик пользователей
	userStore := inmemstorage.NewUserStore()
	userService := user.NewUserService(userStore)
	userHandler := http.NewUserHandler(userService, sessionManager)

	// Создаем http роутер
	r := chi.NewRouter()
	// Register middleware first
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Добавляем swagger
	r.Mount("/swagger", httpSwagger.WrapHandler)

	// Прикрепляем обработчики
	taskHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)

	// Создаем новый сервер с заданной конфигурацией и хранилищем
	srv := httpServer.Server{
		Config: config,
	}

	// Логируем сообщение о запуске сервера
	log.Printf("starting server at %s\n", srv.Config.Addr)
	// Запускаем сервер и проверяем на ошибки
	if err := srv.Run(r); err != nil {
		log.Fatal("server down unexpectedly")
	}
}
