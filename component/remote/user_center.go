package remote

import (
	"github.com/go-resty/resty/v2"
	"github.com/ihezebin/olympus/httpclient"
	"github.com/pkg/errors"
)

var userCenterClient *resty.Client

func InitUserCenter(host string) error {
	userCenterClient = httpclient.NewClient(httpclient.WithHost(host))
	resp, err := userCenterClient.NewRequest().Get("/health")
	if err != nil {
		return errors.Wrap(err, "failed to check user center health")
	}
	if resp.IsError() {
		return errors.New(resp.String())
	}

	return nil
}

func UserCenter() *resty.Client {
	return userCenterClient
}
