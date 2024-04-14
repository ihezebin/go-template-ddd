package pubsub

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

var pulsarClient pulsar.Client

func InitPulsarClient(url string) error {
	var err error
	pulsarClient, err = pulsar.NewClient(pulsar.ClientOptions{
		URL:               url,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		return err
	}

	return nil
}

// PulsarClient https://pulsar.apache.org/docs/3.2.x/client-libraries-go-use/
func PulsarClient() pulsar.Client {
	return pulsarClient
}
