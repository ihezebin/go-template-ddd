package pubsub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
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
