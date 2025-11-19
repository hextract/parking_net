package notification

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
)

type KafkaConnection struct {
	Writer *kafka.Conn
}

func (kc KafkaConnection) Close() error {
	if kc.Writer == nil {
		return errors.New("kafka connection is nil")
	}
	err := kc.Writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close kafka connection: %w", err)
	}
	return nil
}

func NewKafkaConnection(broker string, topic string) (*KafkaConnection, error) {
	if broker == "" {
		return nil, errors.New("kafka broker is required")
	}
	if topic == "" {
		return nil, errors.New("kafka topic is required")
	}

	writer, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}
	return &KafkaConnection{Writer: writer}, nil
}

func NewEnvKafkaConnection() (*KafkaConnection, error) {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	if broker == "" {
		return nil, errors.New("KAFKA_BROKER environment variable is required")
	}
	if topic == "" {
		return nil, errors.New("KAFKA_TOPIC environment variable is required")
	}
	return NewKafkaConnection(broker, topic)
}
