package config

import (
	"code-tasks/pkg/broker"
	"code-tasks/pkg/http"
)

type Config struct {
	HTTPServer http.HTTPConfig
	Rabbit     broker.RabbitConfig
}
