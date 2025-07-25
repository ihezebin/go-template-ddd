package dto

type GenerateCaptchaReq struct {
	Usage string `json:"usage"`
	Key   string `json:"key"`
}

type GenerateCaptchaResp struct {
}

type VerifyCaptchaReq struct {
	Key     string `json:"key"`
	Usage   string `json:"usage"`
	Captcha string `json:"captcha"`
}

type VerifyCaptchaResp struct {
}
