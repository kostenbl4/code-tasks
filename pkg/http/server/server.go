package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var serverShutdownTimeout = 10 * time.Second

type HTTPConfig struct {
	Addr         string        `env:"SERVER_ADDRESS" yaml:"address" env-required:"true"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" yaml:"write_timeout" env-default:"5s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" yaml:"idle_timeout" env-default:"30s"`
}

type Server struct {
	Config HTTPConfig
	server *http.Server
}

func NewServer(cfg HTTPConfig) *Server {
	return &Server{Config: cfg}
}

func (s *Server) Run(ctx context.Context, handler http.Handler) error {

	srv := &http.Server{
		Addr:         s.Config.Addr,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
		IdleTimeout:  s.Config.IdleTimeout,
		Handler:      handler,
	}

	s.server = srv

	errChan := make(chan error, 1)

	slog.Info("starting http server", slog.String("addr", s.Config.Addr))
	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return s.Stop()
	case err := <-errChan:
		return fmt.Errorf("server run error: %w", err)
	}
}

func (s *Server) Stop() error {
	slog.Info("HTTP server shutting down", slog.String("addr", s.server.Addr))
	shutdownCtx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown the server: %w", err)
	}
	return nil
}
