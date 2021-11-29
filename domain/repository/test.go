package repository

import (
	"context"
	"github.com/whereabouts/web-template-ddd/domain/entity"
	"github.com/whereabouts/web-template-ddd/domain/repository/impl/mongo"
)

// TestRepository 只暴露接口, 具体实现在init中初始化, 如使用mongo
type TestRepository interface {
	Register(ctx context.Context, test *entity.Test) error
}

var gTestRepository TestRepository

func GetTestRepository() TestRepository {
	return gTestRepository
}

func InitTestRepository() {
	gTestRepository = mongo.NewTestRepository()
	//gTestRepository = redis.NewTestRepository()
}
