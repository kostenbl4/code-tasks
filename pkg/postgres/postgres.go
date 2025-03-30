package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host     string `env:"DATABASE_HOST" yaml:"host" env-required:"true"`
	Port     int    `env:"DATABASE_PORT" yaml:"port" env-required:"true"`
	User     string `env:"DATABASE_USER" yaml:"user" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" yaml:"password" env-required:"true"`
	DBName   string `env:"DATABASE_NAME" yaml:"db_name" env-required:"true"`
}

func NewPostgresPool(cfg PostgresConfig) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
