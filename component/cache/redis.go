package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var redisCache redis.UniversalClient

func RedisCacheClient() redis.UniversalClient {
	return redisCache
}

func InitRedisCache(ctx context.Context, addr string, password string) error {
	redisCache = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{addr},
		Password: password,
	})
	if err := redisCache.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "redis ping error")
	}

	return nil
}
