package test

import (
	"context"
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

// Repository 只暴露接口, 具体实现在统一初始化 repository.Init 中通过 test.SetRepository 来配置, 如使用mongo
type Repository interface {
	Register(ctx context.Context, test *entity.Test) error
}

var gRepository Repository

func GetRepository() Repository {
	return gRepository
}

func SetRepository(repo Repository) {
	gRepository = repo
}
