package pubsub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ihezebin/soup/pubsub"
	pubsubPulsar "github.com/ihezebin/soup/pubsub/pulsar"
)

func TestPulsar(t *testing.T) {
	if err := InitPulsarClient("pulsar://localhost:6650"); err != nil {
		t.Fatal(err)
	}

	client := PulsarClient()
	defer client.Close()

	ctx := context.Background()
	// consumer
	go func() {
		consumer, err := client.Subscribe(pulsar.ConsumerOptions{
			Topic:            "persistent://public/default/my-topic",
			SubscriptionName: "my-subscription-1",
			Type:             pulsar.Exclusive,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
		defer consumer.Close()
		defer consumer.Unsubscribe()

		for {
			msg, err := consumer.Receive(ctx)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(string(msg.Payload()))
			consumer.Ack(msg)
		}
	}()

	// producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "persistent://public/default/my-topic",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	// send msg
	for i := 0; i < 5; i++ {
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte("hello" + fmt.Sprint(i)),
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("send success", i)

		time.Sleep(time.Second)
	}

	time.Sleep(time.Second * 10)
}

func TestPulsarSubscribe(t *testing.T) {
	if err := InitPulsarSubscribe(pubsubPulsar.SubOptions{
		ClientOptions: pulsar.ClientOptions{
			URL: "pulsar://localhost:6650",
		},
		ConsumerOptions: pulsar.ConsumerOptions{
			Topic:            "persistent://public/default/test",
			SubscriptionName: "test",
			Type:             pulsar.Exclusive,
		},
	}); err != nil {
		t.Fatal(err)
	}

	consumer := PulsarSubscriber()
	go func() {
		time.Sleep(time.Second * 30)
		consumer.Close()
	}()

	ctx := context.Background()

	consumer.Receive(ctx, func(ctx context.Context, msg pubsub.ConsumerMessage) error {
		t.Log(string(msg.Payload()))
		return nil
	})

	t.Log("receive success")
}

func TestPulsarPublish(t *testing.T) {
	if err := InitPulsarPublish(pubsubPulsar.PubOptions{
		ClientOptions: pulsar.ClientOptions{
			URL: "pulsar://localhost:6650",
		},
		ProducerOptions: pulsar.ProducerOptions{
			Topic: "persistent://public/default/test",
		},
	}); err != nil {
		t.Fatal(err)
	}

	publisher := PulsarPublisher()
	ctx := context.Background()
	err := publisher.Send(ctx, pubsub.ProducerMessage{
		Payload: []byte("hello"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send success")
}
