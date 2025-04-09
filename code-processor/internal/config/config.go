package config

import (
	"code-tasks/pkg/broker"
	pkgLogger "code-tasks/pkg/log"
)

type Config struct {
	Rabbit broker.RabbitConfig    `yaml:"rabbit"`
	Logger pkgLogger.LoggerConfig `yaml:"logger"`
}
