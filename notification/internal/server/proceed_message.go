package server

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func GetFormat(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "format" {
			return string(header.Value)
		}
	}
	return ""
}

func (server *NotificationKafkaServer) ProceedMessage(ctx context.Context, key []byte, value []byte, headers []kafka.Header) error {
	if len(value) == 0 {
		return fmt.Errorf("message value is empty")
	}

	format := GetFormat(headers)
	if format != "json" {
		return fmt.Errorf("invalid format: expected 'json', got '%s'", format)
	}

	handlerKey := string(key)
	if handlerKey == "" {
		return fmt.Errorf("handler key is empty")
	}

	handler, exists := server.handlers[handlerKey]
	if !exists {
		return fmt.Errorf("handler not found for key: %s", handlerKey)
	}

	return handler(value)
}
