package service

import (
	"context"

	"github.com/ihezebin/soup/httpserver"
	"github.com/ihezebin/soup/logger"

	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/domain/entity"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/domain/service"
)

type ExampleApplicationService struct {
	logger               *logger.Entry
	exampleRepository    repository.ExampleRepository
	exampleDomainService service.ExampleDomainService
}

func NewExampleApplicationService(l *logger.Entry) *ExampleApplicationService {
	return &ExampleApplicationService{
		logger:               l.WithField("application", "example"),
		exampleRepository:    repository.GetExampleRepository(),
		exampleDomainService: service.GetExampleDomainService(),
	}
}

func (svc *ExampleApplicationService) Login(ctx context.Context, req dto.ExampleLoginReq) (*dto.ExampleLoginResp, error) {
	// application service中通过调用各个 domain 中的 service 或 repository, 实现业务逻辑的编排；
	if ok, errMsg := svc.exampleDomainService.ValidateExample(&entity.Example{
		Username: req.Username,
		Password: req.Password,
	}); !ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, errMsg)
	}

	example, err := svc.exampleRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "find example by username err, example: %+v", example)
		return nil, httpserver.ErrorWithInternalServer()
	}

	if example == nil {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号不存在")
	}

	if !example.CheckPasswordMatch(req.Password) {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "密码不正确")
	}

	token, err := svc.exampleDomainService.GenerateToken(example)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "generate token err, example: %+v", example)
		return nil, httpserver.ErrorWithInternalServer()
	}

	return &dto.ExampleLoginResp{
		Token: token,
	}, nil
}

func (svc *ExampleApplicationService) Register(ctx context.Context, req dto.ExampleRegisterReq) (*dto.ExampleRegisterResp, error) {
	// application service中通过调用各个 domain 中的 service 或 repository, 实现业务逻辑的编排；

	newExample := &entity.Example{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if ok, errMsg := svc.exampleDomainService.ValidateExample(newExample); !ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, errMsg)
	}

	example, err := svc.exampleRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "find example by username err, example: %+v", example)
		return nil, httpserver.ErrorWithInternalServer()
	}
	if example != nil {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号已存在")
	}

	example, err = svc.exampleRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "find example by email err, example: %+v", example)
	}
	if example != nil {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "邮箱已绑定账号")
	}

	newExample.Salt = "xxxx"
	newExample.Password = newExample.MD5PasswordWithSalt()

	if err := svc.exampleRepository.InsertOne(ctx, newExample); err != nil {
		svc.logger.WithError(err).Errorf(ctx, "insert example err, example: %+v", newExample)
		return nil, httpserver.ErrorWithInternalServer()
	}

	// 将敏感信息置空
	newExample = newExample.Sensitive()

	return &dto.ExampleRegisterResp{
		Example: newExample,
	}, nil
}
