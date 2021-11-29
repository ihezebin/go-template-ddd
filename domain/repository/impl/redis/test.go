package redis

import (
	"context"
	"fmt"
	"github.com/whereabouts/sdk/db/redisc"
	"github.com/whereabouts/web-template-ddd/component/cache"
	"github.com/whereabouts/web-template-ddd/domain/entity"
)

type testRepository struct {
	*redisc.Base
}

func NewTestRepository() *testRepository {
	return &testRepository{redisc.NewBaseModel(cache.GetRedisCli(), "captcha")}
}

func (repo *testRepository) Register(ctx context.Context, test *entity.Test) error {
	data, err := test.String()
	if err != nil {
		return err
	}
	err = repo.Set(repo.testKey(test.Name), data).Error()
	return err
}

func (repo *testRepository) testKey(name string) string {
	return fmt.Sprintf("test.name.%s", name)

}
