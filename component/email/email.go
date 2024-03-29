package email

import "github.com/ihezebin/sdk/emailc"

var gClient *emailc.Client

func Init(config emailc.Config) (err error) {
	gClient, err = emailc.NewClientWithConfig(config)
	return
}

func GetCli() *emailc.Client {
	return gClient
}
