package config

import (
	"code-tasks/pkg/broker"
	rediscache "code-tasks/pkg/cache/redis"
	"code-tasks/pkg/http"
	"code-tasks/pkg/postgres"
)

type Config struct {
	HTTPServer http.HTTPConfig         `yaml:"http_server"`
	Rabbit     broker.RabbitConfig     `yaml:"rabbit"`
	Postgres   postgres.PostgresConfig `yaml:"postgres"`
	Redis      rediscache.RedisConfig  `yaml:"redis"`
}
