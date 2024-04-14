package pubsub

import (
	"context"

	"github.com/segmentio/kafka-go"
)

var kafkaConn *kafka.Conn

func InitKafkaConn(ctx context.Context, address, topic string, partition int) error {
	var err error
	kafkaConn, err = kafka.DialLeader(ctx, "tcp", address, topic, partition)
	if err != nil {
		return err
	}

	return nil
}

func KafkaConn() *kafka.Conn {
	return kafkaConn
}
