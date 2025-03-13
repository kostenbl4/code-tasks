package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server - структура сервера
type Server struct {
	Config Config
}

// Config - структура конфигурации сервера
type Config struct {
	Addr string
}

// Run - запуск сервера
func (s *Server) Run(r chi.Router) error {

	srv := http.Server{
		Addr:    s.Config.Addr,
		Handler: r,
	}
	return srv.ListenAndServe()
}
