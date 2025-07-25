package sms

import (
	"context"
	"fmt"

	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/olympus/sms/aliyun"
	"github.com/ihezebin/olympus/sms/tencent"
	"github.com/pkg/errors"
)

type SmsClient interface {
	SendCapatchMessage(ctx context.Context, phone string, captcha string) error
}

type SmsTencentClient struct {
	kernel *tencent.Client
	config config.SmsTencentConfig
}

var gTencentClient *SmsTencentClient
var _ SmsClient = (*SmsTencentClient)(nil)

func (c *SmsTencentClient) SendCapatchMessage(ctx context.Context, phone string, captcha string) error {
	msg := tencent.NewMessage().WithAppId(c.config.AppId).WithSignName(c.config.SignName).
		WithTemplate(c.config.CaptchaTemplateId, captcha, "10")

	faileds, err := c.kernel.Send(ctx, msg, fmt.Sprintf("+86%s", phone))
	if err != nil {
		return errors.Wrapf(err, "send sms msg err, failed reasons: %+v", faileds)
	}

	return nil
}

func ClientTencent() *SmsTencentClient {
	return gTencentClient
}

func InitTencent(conf config.SmsTencentConfig) error {
	client, err := tencent.NewClient(conf.Config)
	if err != nil {
		return err
	}
	gTencentClient = &SmsTencentClient{
		kernel: client,
		config: conf,
	}
	return nil
}

type SmsAliyunClient struct {
	config config.SmsAliyunConfig
	kernel *aliyun.Client
}

var _ SmsClient = (*SmsAliyunClient)(nil)

func (c *SmsAliyunClient) SendCapatchMessage(ctx context.Context, phone string, captcha string) error {
	msg := aliyun.NewMessage().WithSignName(c.config.SignName).WithTemplate(c.config.CaptchaTemplateCode, map[string]interface{}{
		"code": captcha,
	})
	err := c.kernel.Send(ctx, msg, fmt.Sprintf("+86%s", phone))
	if err != nil {
		return errors.Wrap(err, "send sms msg err")
	}
	return nil
}

var gAliyunClient *SmsAliyunClient

func ClientAliyun() *SmsAliyunClient {
	return gAliyunClient
}

func InitAliyun(conf config.SmsAliyunConfig) error {
	client, err := aliyun.NewClient(conf.Config)
	if err != nil {
		return errors.Wrap(err, "init aliyun client error")
	}
	gAliyunClient = &SmsAliyunClient{
		config: conf,
		kernel: client,
	}
	return nil
}
