package router

import (
	"github.com/gin-gonic/gin"
	valication "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/application/service"
	"github.com/ihezebin/go-template-ddd/component/constant"
	"github.com/ihezebin/olympus/httpserver"
	"github.com/ihezebin/olympus/logger"
)

type CaptchaRouter struct {
	logger  logger.Logger
	service *service.CaptchaApplicationService
}

func NewCaptchaRouter() *CaptchaRouter {
	l := logger.WithField("router", "captcha")
	return &CaptchaRouter{
		logger:  l,
		service: service.NewCaptchaApplicationService(l),
	}
}

func (r *CaptchaRouter) RegisterRoutes(router httpserver.Router) {
	group := router.Group("captcha")
	group.POST("/generate", httpserver.NewHandler(r.Generate),
		httpserver.WithOpenAPISummary("生成验证码"),
		httpserver.WithOpenAPIDescription("生成验证码"),
	)
	group.POST("/verify", httpserver.NewHandler(r.Verify),
		httpserver.WithOpenAPISummary("验证验证码"),
		httpserver.WithOpenAPIDescription("验证验证码"),
	)
}

func (r *CaptchaRouter) Generate(c *gin.Context, req dto.GenerateCaptchaReq) (*dto.GenerateCaptchaResp, error) {
	ctx := c.Request.Context()

	usages := make([]interface{}, 0, len(constant.AllowedCaptchaUsages))
	for _, usage := range constant.AllowedCaptchaUsages {
		usages = append(usages, usage)
	}

	if err := valication.ValidateStruct(&req,
		valication.Field(&req.Key, valication.Required),
		valication.Field(&req.Usage, valication.Required, valication.In(usages...)),
	); err != nil {
		r.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return r.service.Generate(ctx, req)
}

func (r *CaptchaRouter) Verify(c *gin.Context, req dto.VerifyCaptchaReq) (*dto.VerifyCaptchaResp, error) {
	ctx := c.Request.Context()

	usages := make([]interface{}, 0, len(constant.AllowedCaptchaUsages))
	for _, usage := range constant.AllowedCaptchaUsages {
		usages = append(usages, usage)
	}

	if err := valication.ValidateStruct(&req,
		valication.Field(&req.Key, valication.Required),
		valication.Field(&req.Usage, valication.Required, valication.In(usages...)),
		valication.Field(&req.Captcha, valication.Required),
	); err != nil {
		r.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return r.service.Verify(ctx, req)
}
