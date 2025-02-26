package service

import (
	"context"

	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type ExampleDomainService interface {
	IsEmailAlreadyExists(ctx context.Context, example *entity.Example) (bool, error)
	IsUsernameAlreadyExists(ctx context.Context, example *entity.Example) (bool, error)
	// 可拆为单独的 token 生成器放于 domain service
	GenerateToken(example *entity.Example) (string, error)
}

var exampleDomainSvc ExampleDomainService

func GetExampleDomainService() ExampleDomainService {
	return exampleDomainSvc
}

func SetExampleDomainService(service ExampleDomainService) {
	exampleDomainSvc = service
}
