package proto

import "github.com/whereabouts/web-template-ddd/domain/entity"

type TestRegisterReq struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

type TestRegisterResp struct {
	Test *entity.Test `json:"welcome"`
}
