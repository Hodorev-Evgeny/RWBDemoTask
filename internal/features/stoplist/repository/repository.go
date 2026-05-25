package feature_repository_stoplist

import core_redis "RWBDwmoTask/internal/core/repository/redis"

type RepositoryStopList struct {
	rds *core_redis.RedisClient
}

func NewRepositoryStopList(
	rds *core_redis.RedisClient,
) *RepositoryStopList {
	return &RepositoryStopList{
		rds: rds,
	}
}
