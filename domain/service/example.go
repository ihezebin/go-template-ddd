package service

import (
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type ExampleDomainService interface {
	ValidateExample(example *entity.Example) (bool, string)
	GenerateToken(example *entity.Example) (string, error)
}

var exampleDomainSvc ExampleDomainService

func GetExampleDomainService() ExampleDomainService {
	return exampleDomainSvc
}

func SetExampleDomainService(service ExampleDomainService) {
	exampleDomainSvc = service
}
