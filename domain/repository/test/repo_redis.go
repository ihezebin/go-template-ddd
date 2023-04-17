package test

import (
	"context"
	"fmt"
	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/ihezebin/sdk/model/redisc"
)

type repoRedis struct {
	redisc.Model
}

func NewRepoRedis(db string) *repoRedis {
	return &repoRedis{Model: redisc.NewModel(cache.GetRedisCli(), "captcha")}
}

func (repo *repoRedis) Register(ctx context.Context, test *entity.Test) error {
	data, err := test.String()
	if err != nil {
		return err
	}
	err = repo.Set(ctx, repo.testKey(test.Name), data, 0).Err()
	return err
}

func (repo *repoRedis) testKey(name string) string {
	return fmt.Sprintf("test.name.%s", name)

}
