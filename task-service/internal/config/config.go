package config

import (
	"github.com/kostenbl4/code-tasks/pkg/broker"
	rediscache "github.com/kostenbl4/code-tasks/pkg/cache/redis"
	"github.com/kostenbl4/code-tasks/pkg/http/server"
	pkgLogger "github.com/kostenbl4/code-tasks/pkg/log"
	"github.com/kostenbl4/code-tasks/pkg/postgres"
)

type Config struct {
	HTTPServer server.HTTPConfig       `yaml:"http_server"`
	Rabbit     broker.RabbitConfig     `yaml:"rabbit"`
	Postgres   postgres.PostgresConfig `yaml:"postgres"`
	Redis      rediscache.RedisConfig  `yaml:"redis"`
	Logger     pkgLogger.LoggerConfig  `yaml:"logger"`
}
