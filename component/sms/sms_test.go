package sms

import (
	"context"
	"testing"

	"github.com/ihezebin/go-template-ddd/config"
	"github.com/ihezebin/olympus/sms/tencent"
)

func TestSms(t *testing.T) {
	err := Init(config.SmsConfig{
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

	err = gClient.SendCapatchMessage(context.Background(), "+8613518468111")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send sms succeed")
}
