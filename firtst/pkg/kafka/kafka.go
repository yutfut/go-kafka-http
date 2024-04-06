package kafka

import (
	"context"
	"fmt"

	"first/pkg/conf"

	"github.com/segmentio/kafka-go"
)

func KafkaConn(config *conf.Conf) (*kafka.Conn, *kafka.Reader, error) {
	address := fmt.Sprintf("%s:%d", config.Kafka.Host, config.Kafka.Port)

	producer, err := kafka.DialLeader(context.Background(), "tcp", address, "input", 0)
	if err != nil {
		return nil, nil, err
	}

	concumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		GroupID:   "consumer-group-id",
		Topic:     "output",
		MaxBytes:  10e6, // 10MB
	})

	return producer, concumer, nil
}