package core_redis

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	ADDR     string `envconfig:"ADDR"`
	Password string `envconfig:"PASSWORD"`
	Database int    `envconfig:"DATABASE"`
}

func getRedisConfig() (RedisConfig, error) {
	redisConfig := RedisConfig{}
	if err := envconfig.Process("REDIS", &redisConfig); err != nil {
		return RedisConfig{}, fmt.Errorf("failed to process redis env var: %w", err)
	}

	return redisConfig, nil
}

func MustGetRedisConfig() RedisConfig {
	redisConfig, err := getRedisConfig()
	if err != nil {
		panic(err)
	}
	return redisConfig
}
