package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type HTTPConfig struct {
	Addr         string        `env:"SERVER_ADDRESS" yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" yaml:"idle_timeout" env-default:"30s"`
}

// Server - структура сервера
type Server struct {
	Config HTTPConfig
}

// Run - запуск сервера
func (s *Server) Run(r chi.Router) error {

	srv := http.Server{
		Addr:         s.Config.Addr,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
		IdleTimeout:  s.Config.IdleTimeout,
		Handler:      r,
	}
	return srv.ListenAndServe()
}
