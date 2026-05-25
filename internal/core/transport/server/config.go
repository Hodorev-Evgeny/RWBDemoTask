package core_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ServerConfig struct {
	Addr            string         `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration  `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
	TimeZone        *time.Location `envconfig:"TIMEZONE" required:"true"`
}

func NewServerConfig() (ServerConfig, error) {
	var cfg ServerConfig

	if err := envconfig.Process("HTTP", &cfg); err != nil {
		return ServerConfig{}, fmt.Errorf("could not process server config: %w", err)
	}

	return cfg, nil
}

func MustNewConfigServer() ServerConfig {
	cfg, err := NewServerConfig()
	if err != nil {
		panic(err)
	}

	return cfg
}
