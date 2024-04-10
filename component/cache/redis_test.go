package cache

import (
	"context"
	"testing"
)

func TestRedis(t *testing.T) {
	ctx := context.Background()
	if err := InitRedisCache(ctx, "127.0.0.1:6379", "root"); err != nil {
		t.Fatal(err)
	}

	client := RedisCacheClient()
	if err := client.Do(ctx, "SET", "key", "value").Err(); err != nil {
		t.Fatal(err)
	}

	val, err := client.Get(ctx, "key").Result()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(val)

}
