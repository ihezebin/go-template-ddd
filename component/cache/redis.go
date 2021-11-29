package cache

import (
	"context"
	"github.com/whereabouts/sdk/db/redisc"
)

var gRedisCli redisc.Client

func InitRedis(ctx context.Context, config redisc.Config) (err error) {
	gRedisCli, err = redisc.NewClientWithConfig(config)
	return
}

func GetRedisCli() redisc.Client {
	return gRedisCli
}
