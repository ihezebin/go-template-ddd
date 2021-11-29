package sms

import smsc "github.com/whereabouts/sdk/smsc/tencent"

var gClient *smsc.Client

func Init(config smsc.Config) (err error) {
	gClient, err = smsc.NewClientWithConfig(config)
	return
}

func GetCli() *smsc.Client {
	return gClient
}
