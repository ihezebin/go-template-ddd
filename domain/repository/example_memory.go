package repository

import (
	"context"

	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type exampleMemoryRepository struct {
	lowLayer ExampleRepository
}

func NewExampleMemoryRepository(lowLayer ExampleRepository) ExampleRepository {
	return &exampleMemoryRepository{
		lowLayer: lowLayer,
	}
}

var _ ExampleRepository = (*exampleMemoryRepository)(nil)

func (e *exampleMemoryRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	return e.lowLayer.InsertOne(ctx, example)
}

func (e *exampleMemoryRepository) FindByUsername(ctx context.Context, username string) (example *entity.Example, err error) {
	key := exampleFindByUsernameKey(username)
	val, ok := cache.GetMemoryCache().Get(key)
	if ok {
		return val.(*entity.Example), nil
	}

	example, err = e.lowLayer.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if example != nil {
		cache.GetMemoryCache().SetDefault(key, example)
		return example, nil
	}

	// 缓存零值
	cache.GetEmptyCache().SetDefault(key, example)

	return nil, nil
}

func (e *exampleMemoryRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	return e.lowLayer.FindByEmail(ctx, email)
}
