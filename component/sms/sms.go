package sms

import (
	"context"
	"fmt"

	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/olympus/sms/tencent"
	"github.com/pkg/errors"
)

type SmsClient struct {
	kernel *tencent.Client
	config config.SmsConfig
}

var gClient *SmsClient

func (c *SmsClient) SendCapatchMessage(ctx context.Context, phone string, params ...string) error {
	msg := tencent.NewMessage().WithAppId(c.config.AppId).WithSignName(c.config.SignName).
		WithTemplate(c.config.CaptchaTemplateId, params...)

	faileds, err := c.kernel.Send(ctx, msg, fmt.Sprintf("+86%s", phone))
	if err != nil {
		return errors.Wrapf(err, "send sms msg err, failed reasons: %+v", faileds)
	}

	return nil
}

func Client() *SmsClient {
	return gClient
}

func Init(conf config.SmsConfig) error {
	client, err := tencent.NewClient(conf.Config)
	if err != nil {
		return err
	}
	gClient = &SmsClient{
		kernel: client,
		config: conf,
	}
	return nil
}
