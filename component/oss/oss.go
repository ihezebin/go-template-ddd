package oss

import (
	"github.com/ihezebin/oneness/oss"
)

var client oss.Client

func Client() oss.Client {
	return client
}

func Init(dsn string) error {
	var err error
	client, err = oss.NewClient(dsn)
	if err != nil {
		return err
	}
	return nil
}
