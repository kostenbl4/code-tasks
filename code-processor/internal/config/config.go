package config

import "code-tasks/pkg/broker"

type Config struct {
	Rabbit broker.RabbitConfig
}
