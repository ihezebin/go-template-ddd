package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ihezebin/go-template-ddd/component/cache"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/pkg/errors"
)

type exampleRedisRepository struct {
	lowLayer ExampleRepository
}

func NewExampleRedisRepository(lowLayer ExampleRepository) ExampleRepository {
	return &exampleRedisRepository{
		lowLayer: lowLayer,
	}
}

var _ ExampleRepository = (*exampleRedisRepository)(nil)

func (e *exampleRedisRepository) InsertOne(ctx context.Context, example *entity.Example) error {
	return e.lowLayer.InsertOne(ctx, example)
}

func (e *exampleRedisRepository) FindByUsername(ctx context.Context, username string) (*entity.Example, error) {
	key := exampleFindByUsernameKey(username)
	cmd := cache.RedisCacheClient().Get(ctx, key)
	data, err := cmd.Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if data != nil {
		example := &entity.Example{}
		err = cmd.Scan(example)
		if err != nil {
			return nil, err
		}
		return example, nil
	}

	example, err := e.lowLayer.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if example != nil {
		if err = cache.RedisCacheClient().Set(ctx, key, example, time.Minute*30).Err(); err != nil {
			return nil, err
		}
		return example, nil
	}

	return nil, nil
}

func (e *exampleRedisRepository) FindByEmail(ctx context.Context, email string) (example *entity.Example, err error) {
	return e.lowLayer.FindByEmail(ctx, email)
}
