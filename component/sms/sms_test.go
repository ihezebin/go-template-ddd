package sms

import (
	"context"
	"testing"

	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/olympus/logger"
	"github.com/ihezebin/olympus/sms/aliyun"
	"github.com/ihezebin/olympus/sms/tencent"
)

func TestTencentSms(t *testing.T) {
	conf, err := config.Load("../../config/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	logger.Debugf(ctx, "load config: %+v", conf.String())
	err = InitTencent(config.SmsTencentConfig{
		Config: tencent.Config{
			SecretId:  "SecretId",
			SecretKey: "SecretKey",
			Region:    "ap-guangzhou",
		},
		CaptchaTemplateId: "",
		AppId:             "",
		SignName:          "",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ClientTencent().SendCapatchMessage(context.Background(), "13518468111", "123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send sms succeed")
}

func TestAliyunSms(t *testing.T) {
	conf, err := config.Load("../../config/config.toml")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	logger.Debugf(ctx, "load config: %+v", conf.String())
	err = InitAliyun(config.SmsAliyunConfig{
		Config: aliyun.Config{
			AccessKeyId:     "AccessKeyId",
			AccessKeySecret: "AccessKeySecret",
			Endpoint:        "Endpoint",
		},
		SignName:            "SignName",
		CaptchaTemplateCode: "CaptchaTemplateCode",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ClientAliyun().SendCapatchMessage(context.Background(), "13518468111", "123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send sms succeed")
}
