package config

import (
	"github.com/kostenbl4/code-tasks/pkg/broker"
	pkgLogger "github.com/kostenbl4/code-tasks/pkg/log"
)

type Config struct {
	Rabbit broker.RabbitConfig    `yaml:"rabbit"`
	Logger pkgLogger.LoggerConfig `yaml:"logger"`
}
