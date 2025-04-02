package main

import (
	"code-tasks/pkg/broker"
	httpServer "code-tasks/pkg/http"
	"code-tasks/pkg/postgres"
	"code-tasks/task-service/internal/api/http"
	"code-tasks/task-service/internal/config"
	"code-tasks/task-service/internal/middleware/metrics"

	//inmemstorage "code-tasks/task-service/internal/repository/in-mem-storage"
	rediscache "code-tasks/pkg/cache/redis"
	postgresstorage "code-tasks/task-service/internal/repository/postgres_storage"
	rabbitconsumer "code-tasks/task-service/internal/repository/rabbit_consumer"
	rabbitsender "code-tasks/task-service/internal/repository/rabbit_sender"
	redisstorage "code-tasks/task-service/internal/repository/redis_storage"
	"code-tasks/task-service/internal/usecases/session"
	"code-tasks/task-service/internal/usecases/task"
	"code-tasks/task-service/internal/usecases/user"
	"log"

	_ "code-tasks/task-service/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyakaznacheev/cleanenv"

	//"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title task-service API
// @version 1.0
// @description This is a sample server.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {

	var cfg config.Config

	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Fatal(err)
	}

	// Создаем хранилище в операционной памяти
	// sessionStore := inmemstorage.NewSessionStore()

	redis, err := rediscache.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем хранилище redis
	sessionStore := redisstorage.NewSessionStore(redis, cfg.Redis.TTL)

	// Cоздаем менеджер сессий
	sessionManager := session.NewSeessionManager(sessionStore, 3600)

	// Создаем хранилище postgres
	PGpool, err := postgres.NewPostgresPool(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer PGpool.Close()

	// Создаем хранилище, сервис, обработчик задач
	taskStore := postgresstorage.NewTaskStore(PGpool)

	sendConn, err := broker.ConnectRabbitMQ(cfg.Rabbit)
	if err != nil {
		log.Fatal(err)
	}
	sendClient, err := broker.NewRabbitClient(sendConn)
	if err != nil {
		log.Fatal(err)
	}
	taskSender := rabbitsender.New(sendClient)

	consumeConn, err := broker.ConnectRabbitMQ(cfg.Rabbit)
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

	// go taskService.ListenTaskProcessor()

	// Создаем хранилище, сервис, обработчик пользователей
	userStore := postgresstorage.NewUserStore(PGpool)
	userService := user.NewUserService(userStore)
	userHandler := http.NewUserHandler(userService, sessionManager)

	// Создаем http роутер
	r := chi.NewRouter()
	// Register middleware first
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	

	// Создаем отдельный Registry для метрик
	registry := prometheus.NewRegistry()
	
	// Добавляем go runtime metrics и process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	r.Use(metrics.New(registry, nil, "task_service").Handler())

	r.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Добавляем swagger
	r.Mount("/swagger", httpSwagger.WrapHandler)

	// Прикрепляем обработчики
	taskHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)

	// Создаем новый сервер с заданной конфигурацией и хранилищем
	srv := httpServer.Server{
		Config: cfg.HTTPServer,
	}

	// Логируем сообщение о запуске сервера
	log.Printf("starting server at %s\n", srv.Config.Addr)
	// Запускаем сервер и проверяем на ошибки
	if err := srv.Run(r); err != nil {
		log.Fatal("server down unexpectedly")
	}

	// TODO: graceful shutdown, soon)
}
