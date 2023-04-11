package test

import "github.com/ihezebin/web-template-ddd/domain/entity"

type RegisterReq struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

type RegisterResp struct {
	Test *entity.Test `json:"welcome"`
}
