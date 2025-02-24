package emailc

import (
	"github.com/ihezebin/soup/email"
)

var client *email.Client

func Client() *email.Client {
	return client
}

func Init(conf email.Config) error {
	var err error
	client, err = email.NewClient(conf)
	if err != nil {
		return err
	}
	return nil
}
