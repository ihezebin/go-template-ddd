package dto

import "github.com/ihezebin/go-template-ddd/component/constant"

type GenerateCaptchaReq struct {
	Usage constant.CaptchaUsageType `json:"usage"`
	Key   string                    `json:"key"`
}

type GenerateCaptchaResp struct {
}

type VerifyCaptchaReq struct {
	Key     string                    `json:"key"`
	Usage   constant.CaptchaUsageType `json:"usage"`
	Captcha string                    `json:"captcha"`
}

type VerifyCaptchaResp struct {
}
