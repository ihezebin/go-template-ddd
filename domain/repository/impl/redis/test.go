package redis

import (
	"context"
	"fmt"
	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/ihezebin/sdk/model/redisc"
)

type testRepository struct {
	redisc.Model
}

func NewTestRepository() *testRepository {
	return &testRepository{Model: redisc.NewModel(cache.GetRedisCli(), "captcha")}
}

func (repo *testRepository) Register(ctx context.Context, test *entity.Test) error {
	data, err := test.String()
	if err != nil {
		return err
	}
	err = repo.Set(ctx, repo.testKey(test.Name), data, 0).Err()
	return err
}

func (repo *testRepository) testKey(name string) string {
	return fmt.Sprintf("test.name.%s", name)

}
