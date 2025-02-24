package pubsub

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ihezebin/soup/pubsub"
	pubsubPulsar "github.com/ihezebin/soup/pubsub/pulsar"
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

var pulsarPublisher pubsub.Publisher

func PulsarPublisher() pubsub.Publisher {
	return pulsarPublisher
}

func InitPulsarPublish(option pubsubPulsar.PubOptions) error {
	var err error
	pulsarPublisher, err = pubsubPulsar.NewPublisher(option)
	if err != nil {
		return err
	}

	return nil
}

var pulsarSubscriber pubsub.Subscriber

func PulsarSubscriber() pubsub.Subscriber {
	return pulsarSubscriber
}

func InitPulsarSubscribe(option pubsubPulsar.SubOptions) error {
	var err error
	pulsarSubscriber, err = pubsubPulsar.NewSubscriber(option)
	if err != nil {
		return err
	}

	return nil
}
