package service

import (
	"github.com/ihezebin/go-template-ddd/domain/entity"
)

type ExampleDomainService interface {
	ValidateExample(example *entity.Example) (bool, string)
	GenerateToken(example *entity.Example) (string, error)
}

var exampleDomainService ExampleDomainService

func GetExampleDomainService() ExampleDomainService {
	return exampleDomainService
}

func SetExampleDomainService(service ExampleDomainService) {
	exampleDomainService = service
}
