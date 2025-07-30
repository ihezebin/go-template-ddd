package constant

const (
	TokenSecret = "665590964a0cb0c58bb5fb3b"

	HeaderKeyToken = "Token"
	HeaderKeyUid   = "H-Uid"

	QueryKeyUid = "h_uid"

	UsernameAdmin = "hezebin"
)

type CaptchaUsageType string

func (c CaptchaUsageType) String() string {
	return string(c)
}

const (
	LoginCaptchaUsage    CaptchaUsageType = "login"
	RegisterCaptchaUsage CaptchaUsageType = "register"
)

var AllowedCaptchaUsages = []CaptchaUsageType{
	LoginCaptchaUsage,
	RegisterCaptchaUsage,
}
