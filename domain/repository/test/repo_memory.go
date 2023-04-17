package test

import (
	"context"
	"fmt"
	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	gocache "github.com/patrickmn/go-cache"
)

type repoMemory struct {
	memoryCache *gocache.Cache
	emptyCache  *gocache.Cache
}

func NewRepoMemory(db string) *repoMemory {
	return &repoMemory{
		memoryCache: cache.GetMemoryCache(),
		emptyCache:  cache.GetEmptyCache(),
	}
}

func (repo *repoMemory) Register(ctx context.Context, test *entity.Test) error {
	data, err := test.String()
	if err != nil {
		return err
	}
	repo.memoryCache.Set(repo.testKey(test.Name), data, -1)
	return nil
}

func (repo *repoMemory) testKey(name string) string {
	return fmt.Sprintf("test.name.%s", name)

}
