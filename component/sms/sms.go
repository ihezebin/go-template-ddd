package sms

import smsc "github.com/ihezebin/sdk/smsc/tencent"

var gClient *smsc.Client
var gMsg *smsc.Message

func Init(config smsc.Config, message smsc.Message) (err error) {
	gClient, err = smsc.NewClientWithConfig(config)
	gMsg = &message
	return
}

func GetMessage() *smsc.Message {
	return gMsg
}

func GetCli() *smsc.Client {
	return gClient
}
