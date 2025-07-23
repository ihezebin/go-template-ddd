package sms

import (
	"context"
	"testing"

	"github.com/ihezebin/olympus/sms/tencent"
)

func TestEmail(t *testing.T) {
	err := Init(tencent.Config{
		SecretId:  "SecretId",
		SecretKey: "SecretKey",
		Region:    "ap-guangzhou",
	})
	if err != nil {
		t.Fatal(err)
	}

	msg := tencent.NewMessage().WithAppId("1400578890").WithSignName("hezebin").
		WithTemplate("11477481", 123321, 10)
	faileds, err := client.Send(context.Background(), msg, "+8613518468111")
	if err != nil {
		t.Error(faileds)
		t.Fatal(err)
	}
	t.Log("send sms succeed")
}
