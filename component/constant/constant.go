package constant

const (
	TokenSecret = "665590964a0cb0c58bb5fb3b"

	HeaderKeyToken = "Token"
	HeaderKeyUid   = "H-Uid"

	QueryKeyUid = "h_uid"

	UsernameAdmin = "hezebin"
)

const (
	LoginCaptchaUsage    = "login"
	RegisterCaptchaUsage = "register"
)

var AllowedCaptchaUsages = []string{
	LoginCaptchaUsage,
	RegisterCaptchaUsage,
}
