package service

import (
	"context"
	"testing"

	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/domain/repository"
	"github.com/ihezebin/go-template-ddd/domain/service"
)

func TestExampleService(t *testing.T) {
	exampleRepository := repository.NewExampleMockRepository()

	svc := &ExampleApplicationService{
		exampleDomainService: service.NewExampleServiceMock(exampleRepository),
		exampleRepository:    exampleRepository,
		passwordEncoder:      service.NewMd5WithSaltPasswordEncoder(),
	}

	ctx := context.Background()

	registerReq := dto.ExampleRegisterReq{Username: "hezebin", Password: "123123123", Email: "hezebin@qq.com"}
	registerResp, err := svc.Register(ctx, registerReq)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", registerResp.Example)

	loginReq := dto.ExampleLoginReq{Username: "hezebin", Password: "123123123"}
	loginResp, err := svc.Login(ctx, loginReq)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(loginResp.Token)
}
