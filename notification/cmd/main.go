package main

import (
	"github.com/h4x4d/parking_net/notification/internal/handlers"
	"github.com/h4x4d/parking_net/notification/internal/server"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"log"
	"log/slog"
	"os"
)

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		log.Fatalln("KAFKA_BROKER environment variable is not set")
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatalln("KAFKA_TOPIC environment variable is not set")
	}

	groupID := os.Getenv("KAFKA_GROUP_ID")
	if groupID == "" {
		log.Fatalln("KAFKA_GROUP_ID environment variable is not set")
	}

	tracer, err := jaeger.InitTracer("Notification")
	if err != nil {
		log.Printf("Warning: failed to initialize tracer: %v. Continuing without tracing.", err)
		tracer = nil
	}

	notifyHandlers := map[string]func([]byte) error{
		"send_notification": handlers.SendNotificationHandler,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	notificationServer := server.NewNotificationServer(&[]string{broker}, &topic, &groupID, notifyHandlers, tracer)

	if err := notificationServer.Serve(); err != nil {
		log.Fatalln(err)
	}
}
