package main

import (
	"flag"
	"log"
	"task-server/internal/api/http"
	inmemstorage "task-server/internal/repository/in-mem-storage"
	service "task-server/internal/usecases/task"
	httpServer "task-server/pkg/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "task-server/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// @title Task-server API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /
func main() {

	addr := flag.String("addr", ":8080", "address for the http server")

	// Создаем конфигурацию для API сервера
	config := httpServer.Config{
		Addr: *addr,
	}

	// Создаем новое хранилище задач
	repo := inmemstorage.NewTaskStore()

	// Создаем новый объект сервис - логики задач
	service := service.NewTaskService(repo)

	// Создаем новый объект api обработчиков задач
	handler := http.NewTaskHandler(service)

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
	handler.RegisterRoutes(r)

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
