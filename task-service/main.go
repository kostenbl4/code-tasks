package main

import (
	"code-tasks/pkg/broker"
	httpLogger "code-tasks/pkg/http/middleware"
	httpServer "code-tasks/pkg/http/server"
	pkgLogger "code-tasks/pkg/log"
	"code-tasks/pkg/postgres"
	"code-tasks/task-service/internal/api/http"
	"code-tasks/task-service/internal/config"
	"code-tasks/task-service/internal/middleware/metrics"
	"context"
	"os"
	"os/signal"
	"syscall"

	//inmemstorage "code-tasks/task-service/internal/repository/in-mem-storage"
	rediscache "code-tasks/pkg/cache/redis"
	postgresstorage "code-tasks/task-service/internal/repository/postgres_storage"

	//rabbitconsumer "code-tasks/task-service/internal/repository/rabbit_consumer"
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

	logger, file := pkgLogger.NewLogger(cfg.Logger)
	defer file.Close()

	// Создаем хранилище в операционной памяти
	// sessionStore := inmemstorage.NewSessionStore()

	redis, err := rediscache.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	sessionStore := redisstorage.NewSessionStore(redis, cfg.Redis.TTL)
	sessionManager := session.NewSeessionManager(sessionStore)

	PGpool, err := postgres.NewPostgresPool(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer PGpool.Close()

	taskStore := postgresstorage.NewTaskStore(PGpool)

	sendConn, err := broker.ConnectRabbitMQ(cfg.Rabbit)
	if err != nil {
		log.Fatal(err)
	}
	defer sendConn.Close()

	sendClient, err := broker.NewRabbitClient(sendConn)
	if err != nil {
		log.Fatal(err)
	}
	taskSender := rabbitsender.New(sendClient)

	// Настройка принятия результатов задач через брокер
	// consumeConn, err := broker.ConnectRabbitMQ(cfg.Rabbit)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer consumeConn.Close()
	// consumeClient, err := broker.NewRabbitClient(consumeConn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// taskConsumer := rabbitconsumer.New(consumeClient)

	taskService := task.NewTaskService(taskStore, taskSender, nil)
	taskHandler := http.NewTaskHandler(logger, taskService, sessionManager)

	// go taskService.ListenTaskProcessor()

	userStore := postgresstorage.NewUserStore(PGpool)
	userService := user.NewUserService(userStore)
	userHandler := http.NewUserHandler(logger, userService, sessionManager)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(httpLogger.NewLoggingMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	r.Use(metrics.New(registry, nil, "task_service").Handler())

	r.Mount("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	r.Mount("/swagger", httpSwagger.WrapHandler)

	taskHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)

	srv := httpServer.NewServer(cfg.HTTPServer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	if err := srv.Run(ctx, r); err != nil {
		log.Fatal("server down unexpectedly")
	}
}
