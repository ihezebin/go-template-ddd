package service

import (
	"context"

	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/domain/service"
	"github.com/ihezebin/olympus/httpserver"
	"github.com/ihezebin/olympus/logger"
)

type CaptchaApplicationService struct {
	logger           logger.Logger
	captchaGenerater service.CaptchaGenerater
}

func NewCaptchaApplicationService(logger logger.Logger) *CaptchaApplicationService {
	return &CaptchaApplicationService{
		logger:           logger.WithField("application", "captcha"),
		captchaGenerater: service.GetCaptchaGenerater(),
	}
}

func (svc *CaptchaApplicationService) Generate(ctx context.Context, req dto.GenerateCaptchaReq) (*dto.GenerateCaptchaResp, error) {
	ok, captcha, frequency, err := svc.captchaGenerater.Generate(ctx, req.Key, req.Usage)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "generate captcha err, key: %s, usage: %s, captcha: %s, frequency: %d", req.Key, req.Usage, captcha, frequency)
		return nil, httpserver.ErrorWithInternalServer()
	}

	if !ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "生成验证码过于频繁")
	}
	return &dto.GenerateCaptchaResp{}, nil
}

func (svc *CaptchaApplicationService) Verify(ctx context.Context, req dto.VerifyCaptchaReq) (*dto.VerifyCaptchaResp, error) {
	ok, err := svc.captchaGenerater.Verify(ctx, req.Key, req.Usage, req.Captcha)
	if err != nil {
		svc.logger.WithError(err).Errorf(ctx, "verify captcha err, key: %s, usage: %s, captcha: %s", req.Key, req.Usage, req.Captcha)
		return nil, httpserver.ErrorWithInternalServer()
	}

	if !ok {
		return nil, httpserver.NewError(httpserver.CodeValidateRuleFailed, "验证码错误")
	}
	return &dto.VerifyCaptchaResp{}, nil
}
