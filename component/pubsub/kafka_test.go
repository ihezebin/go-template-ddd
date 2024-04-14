package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

func TestKafka(t *testing.T) {
	ctx := context.Background()
	err := InitKafkaConn(ctx, "localhost:9092", "quickstart-events", 0)
	if err != nil {
		t.Fatal(err)
	}

	conn := KafkaConn()
	defer conn.Close()

	go func() {
		// 读取一批消息，得到的batch是一系列消息的迭代器
		batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
		defer batch.Close()

		for {
			msg, err := batch.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}

			t.Log(string(msg.Value))
		}
	}()

	// 发送消息
	var n int
	n, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)

	time.Sleep(time.Second * 10)
}

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	err := InitKafkaConn(ctx, "localhost:9092", "quickstart-events", 0)
	if err != nil {
		t.Fatal(err)
	}

	conn := KafkaConn()
	defer conn.Close()

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "",
		Topic:       "quickstart-events",
		MinBytes:    10e3,             // 10KB
		MaxBytes:    10e6,             // 10MB
		StartOffset: kafka.LastOffset, // 这个很关键，决定了是否是从最新的位置消费数据
	})

	go func() {
		// 消费数据
		for {
			msg, err := consumer.ReadMessage(ctx)
			if err != nil {
				t.Fatal(err)
			}

			t.Log(string(msg.Value))
		}
	}()

	time.Sleep(time.Second * 20)
	err = consumer.Close()
	if err != nil {
		t.Fatal(err)
	}
}
