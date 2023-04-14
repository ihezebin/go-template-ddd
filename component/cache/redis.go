package cache

import (
	"context"
	"github.com/ihezebin/sdk/model/redisc"
)

var gRedisCli *redisc.Client

func InitRedis(ctx context.Context, config redisc.Config) (err error) {
	gRedisCli, err = redisc.NewClient(ctx, config)
	return
}

func GetRedisCli() *redisc.Client {
	return gRedisCli
}
