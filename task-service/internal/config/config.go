package config

import (
	"code-tasks/pkg/broker"
	rediscache "code-tasks/pkg/cache/redis"
	"code-tasks/pkg/http/server"
	pkgLogger "code-tasks/pkg/log"
	"code-tasks/pkg/postgres"
)

type Config struct {
	HTTPServer server.HTTPConfig       `yaml:"http_server"`
	Rabbit     broker.RabbitConfig     `yaml:"rabbit"`
	Postgres   postgres.PostgresConfig `yaml:"postgres"`
	Redis      rediscache.RedisConfig  `yaml:"redis"`
	Logger     pkgLogger.LoggerConfig  `yaml:"logger"`
}
