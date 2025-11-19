package server

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type KafkaServer interface {
	Serve() error
	ProceedMessage(ctx context.Context, key []byte, value []byte, headers []kafka.Header) error
}

type NotificationKafkaServer struct {
	Brokers *[]string
	Topic   *string
	GroupID *string

	handlers map[string]func([]byte) error
	reader   *kafka.Reader
	tracer   trace.Tracer

	KafkaServer
}

func NewNotificationServer(brokers *[]string, topic *string, groupID *string,
	handlers map[string]func([]byte) error, tracer trace.Tracer) *NotificationKafkaServer {
	return &NotificationKafkaServer{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: *brokers,
			Topic:   *topic,
			GroupID: *groupID,
		}),
		handlers: handlers,
		tracer:  tracer,
	}
}

func (server *NotificationKafkaServer) Serve() error {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stopChan
		slog.Info("shutdown signal received, exiting gracefully")
		cancel()
	}()

	reader := server.reader
	defer reader.Close()
	for {
		select {
		case <-ctx.Done():
			slog.Info("consumer stopped")
			return nil
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}
				slog.Error("failed to read message from Kafka",
					slog.String("error", err.Error()),
					slog.String("topic", *server.Topic),
				)
				continue
			}

			var span trace.Span
			if server.tracer != nil {
				ctx, span = server.tracer.Start(ctx, "process_kafka_message")
			}

			proceedErr := server.ProceedMessage(ctx, message.Key, message.Value, message.Headers)
			if proceedErr != nil {
				traceID := ""
				if span != nil {
					traceID = fmt.Sprintf("%s", span.SpanContext().TraceID())
					span.End()
				}
				slog.Error("failed to process message",
					slog.String("trace_id", traceID),
					slog.String("error", proceedErr.Error()),
					slog.String("key", string(message.Key)),
					slog.String("topic", *server.Topic),
					slog.Int("partition", message.Partition),
					slog.Int64("offset", message.Offset),
				)
				continue
			}

			if span != nil {
				span.End()
			}
		}
	}
}
