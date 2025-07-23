package sms

import (
	"github.com/ihezebin/olympus/sms/tencent"
)

var client *tencent.Client

func Client() *tencent.Client {
	return client
}

func Init(conf tencent.Config) error {
	var err error
	client, err = tencent.NewClient(conf)
	if err != nil {
		return err
	}
	return nil
}
