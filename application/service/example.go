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
	exampleEsRepository  repository.ExampleRepository
	exampleDomainService service.ExampleDomainService
	passwordEncoder      service.PasswordEncoder
}

func NewExampleApplicationService(l *logger.Entry) *ExampleApplicationService {
	return &ExampleApplicationService{
		logger:               l.WithField("application", "example"),
		exampleRepository:    repository.GetExampleRepository(),
		exampleEsRepository:  repository.GetExampleEsRepository(),
		exampleDomainService: service.GetExampleDomainService(),
		passwordEncoder:      service.NewMd5WithSaltPasswordEncoder(),
	}
}

func (svc *ExampleApplicationService) Login(ctx context.Context, req dto.ExampleLoginReq) (*dto.ExampleLoginResp, error) {
	// application service中通过调用各个 domain 中的 service 或 repository, 实现业务逻辑的编排；
	example := &entity.Example{
		Username: req.Username,
		Password: req.Password,
	}

	if !example.ValidateUsernameRule() {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号格式不正确")
	}
	if !example.ValidatePasswordRule() {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "密码格式不正确")
	}

	example, err := svc.exampleRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "find example by username err, example: %+v", example)
		return nil, httpserver.ErrorWithInternalServer()
	}

	if example == nil {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号不存在")
	}

	ok, err := svc.passwordEncoder.Verify(req.Password, example.Salt, example.Password)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "verify password err, example: %+v", example)
		return nil, httpserver.ErrorWithInternalServer()
	}
	if !ok {
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

	if !newExample.ValidateUsernameRule() {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号格式不正确")
	}
	if !newExample.ValidatePasswordRule() {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "密码格式不正确")
	}
	if !newExample.ValidateEmailRule() {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "邮箱格式不正确")
	}

	ok, err := svc.exampleDomainService.IsEmailAlreadyExists(ctx, newExample)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "is email already exists err, example: %+v", newExample)
		return nil, httpserver.ErrorWithInternalServer()
	}
	if ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "邮箱已绑定账号")
	}

	ok, err = svc.exampleDomainService.IsUsernameAlreadyExists(ctx, newExample)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "is username already exists err, example: %+v", newExample)
		return nil, httpserver.ErrorWithInternalServer()
	}
	if ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "账号已存在")
	}

	newExample.Salt = "xxxx"
	newExample.Password, err = svc.passwordEncoder.Encode(newExample.Password, newExample.Salt)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "encode password err, example: %+v", newExample)
		return nil, httpserver.ErrorWithInternalServer()
	}

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
