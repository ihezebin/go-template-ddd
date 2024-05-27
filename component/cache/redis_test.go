package cache

import (
	"context"
	"testing"
	"time"

	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	if err := InitRedisCache(ctx, "127.0.0.1:6379", "root"); err != nil {
		t.Fatal(err)
	}

	client := RedisCacheClient()

	example := &entity.Example{
		Username: "hezebin",
		Email:    "ihezebin@qq.com",
	}

	if err := client.Set(ctx, "key", example, time.Minute*5).Err(); err != nil {
		t.Fatal(err)
	}

	example = &entity.Example{}
	err := client.Get(ctx, "key").Scan(example)
	if err != nil && !errors.Is(err, redis.Nil) {
		t.Fatal(err)
	}

	t.Log(example)

}
