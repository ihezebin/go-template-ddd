package service

import "github.com/ihezebin/go-template-ddd/domain/service/impl"

func Init() {
	SetExampleDomainService(impl.NewExampleService())
}
