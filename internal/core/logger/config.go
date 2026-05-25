package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Create setting for logger
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("failed to process env config: %w", err)
	}

	return config, nil
}

func MustNewConfig() Config {
	config, err := NewConfig()

	if err != nil {
		err := fmt.Errorf("failed to create config: %w", err)
		panic(err)
	}

	return config
}
