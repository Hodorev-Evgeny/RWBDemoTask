package core_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func newRedisClient(config RedisConfig) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.ADDR,
		Password: config.Password,
		DB:       config.Database,
	})

	return &RedisClient{
		rdb,
	}, nil
}

func CreateRedisClientMust(config RedisConfig) *RedisClient {
	client, err := newRedisClient(config)
	if err != nil {
		panic(err)
	}

	return client
}

func (r *RedisClient) Protect(
	ctx context.Context,
	userID int64,
	session string,
	query string,
) (bool, error) {
	key := fmt.Sprintf("%s:%s:%s", userID, session, query)

	ok, err := r.SetNX(ctx, key, userID, 30*time.Second).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
