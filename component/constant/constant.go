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
	LoginCaptchaUsage          CaptchaUsageType = "login"
	RegisterCaptchaUsage       CaptchaUsageType = "register"
	ResetPasswordCaptchaUsage  CaptchaUsageType = "reset_password"
	ModifyPasswordCaptchaUsage CaptchaUsageType = "modify_password"
)

var AllowedCaptchaUsages = []CaptchaUsageType{
	LoginCaptchaUsage,
	RegisterCaptchaUsage,
	ResetPasswordCaptchaUsage,
	ModifyPasswordCaptchaUsage,
}

type SortType string

func (s SortType) String() string {
	return string(s)
}

const (
	SortTypeAsc       SortType = "asc"
	SortTypeDesc      SortType = "desc"
	SortTypeUndefined SortType = ""
)
