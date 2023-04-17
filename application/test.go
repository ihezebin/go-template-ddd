package application

import (
	"context"
	testRepo "github.com/ihezebin/go-template-ddd/domain/repository/test"
	"github.com/ihezebin/go-template-ddd/domain/service"
	dto "github.com/ihezebin/go-template-ddd/server/dto/test"
	"github.com/ihezebin/sdk/logger"
	"github.com/pkg/errors"
)

type TestApplication struct {
	logger         *logger.Entry
	testRepository testRepo.Repository
	testService    *service.TestService
}

func NewTestApplication(l *logger.Entry) *TestApplication {
	return &TestApplication{
		logger:         l.WithField("application", "test"),
		testRepository: testRepo.GetRepository(),
		testService:    service.NewTestService(),
	}
}

func (app *TestApplication) TestRegister(ctx context.Context, req *dto.RegisterReq) (*dto.RegisterResp, error) {
	// application中通过调用各个service, 实现业务逻辑的编排；
	// 若业务逻辑过于简单（如简单的增删查改），可越过domain层，在application中直接使用repository进行操作

	test, err := app.testService.TestRegister(ctx, req.Name, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "test post err")
	}

	return &dto.RegisterResp{
		Test: test,
	}, nil
}
