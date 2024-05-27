package cache

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var redisCache redis.UniversalClient

func RedisCacheClient() redis.UniversalClient {
	return redisCache
}

func InitRedisCache(ctx context.Context, addrs []string, password string) error {
	redisCache = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: password,
	})
	if err := redisCache.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "redis ping error")
	}

	return nil
}
