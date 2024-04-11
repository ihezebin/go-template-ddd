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
	}

	ctx := context.Background()

	registerResp, err := svc.Register(ctx, &dto.ExampleRegisterReq{Username: "hezebin", Password: "123123123", Email: "hezebin@qq.com"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", registerResp.Example)

	loginResp, err := svc.Login(ctx, &dto.ExampleLoginReq{Username: "hezebin", Password: "123123123"})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(loginResp.Token)
}
