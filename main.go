package main

import (
	"log"
	"task-server/internal/api"
	"task-server/internal/storage"
)

func main() {
	// Создаем конфигурацию для API сервера
	config := api.Config{
		Addr: "localhost:8080",
	}

	// Создаем новое хранилище задач
	store := storage.NewTaskStore()

	// Создаем новый сервер с заданной конфигурацией и хранилищем
	srv := api.Server{
		Config: config,
		Store:  store,
	}

	// Логируем сообщение о запуске сервера
	log.Printf("starting server at %s\n", srv.Config.Addr)
	// Запускаем сервер и проверяем на ошибки
	if err := srv.Run(); err != nil {
		log.Fatal("server down due to error")
	}
}
