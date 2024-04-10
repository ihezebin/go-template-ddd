package dto

import "github.com/ihezebin/go-template-ddd/domain/entity"

type ExampleLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ExampleLoginResp struct {
	Token string `json:"token"`
}

type ExampleRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ExampleRegisterResp struct {
	Example *entity.Example `json:"example"`
}
