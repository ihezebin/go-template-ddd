package application

import (
	"context"
	"github.com/pkg/errors"
	"github.com/whereabouts/web-template-ddd/domain/entity"
	"github.com/whereabouts/web-template-ddd/domain/repository"
	"github.com/whereabouts/web-template-ddd/domain/service"
)

type TestApplication struct {
	testRepository repository.TestRepository
	testService    *service.TestService
}

func NewTestApplication() *TestApplication {
	return &TestApplication{
		testRepository: repository.GetTestRepository(),
		testService:    service.NewTestService(),
	}
}

func (app *TestApplication) TestRegister(ctx context.Context, name string, password string) (*entity.Test, error) {
	// application中通过调用各个service, 实现业务逻辑的编排；
	// 若业务逻辑过于简单（如简单的增删查改），可越过domain层，在application中直接使用repository进行操作

	test, err := app.testService.TestRegister(ctx, name, password)
	if err != nil {
		return nil, errors.Wrap(err, "test post err")
	}

	return test, nil
}
