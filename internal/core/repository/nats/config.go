package core_nats

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type NatsConfig struct {
	URL string `envconfig:"URL"`
}

func NewNatsConfig() (NatsConfig, error) {
	var config NatsConfig
	err := envconfig.Process("NATS", &config)
	if err != nil {
		return NatsConfig{}, fmt.Errorf("could not process NATS config: %w", err)
	}
	return config, nil
}

func MustNewNatsConfig() NatsConfig {
	nc, err := NewNatsConfig()
	if err != nil {
		panic(err)
	}
	return nc
}
