package handlers

import (
	"github.com/h4x4d/parking_net/booking/internal/database_service"
	"github.com/h4x4d/parking_net/pkg/client"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"github.com/h4x4d/parking_net/pkg/notification"
	"go.opentelemetry.io/otel/trace"
	"log"
	"log/slog"
)

type Handler struct {
	Database  *database_service.DatabaseService
	KafkaConn *notification.KafkaConnection
	KeyCloak  *client.Client
	tracer    trace.Tracer
}

func NewHandler(connStr string) (*Handler, error) {
	db, err := database_service.NewDatabaseService(connStr)
	if err != nil {
		return nil, err
	}
	conn, kafkaErr := notification.NewEnvKafkaConnection()
	if kafkaErr != nil {
		log.Printf("Warning: failed to initialize Kafka connection: %v. Continuing without notifications.", kafkaErr)
		conn = nil
	}
	keycloakClient, keycloakErr := client.NewClient()
	if keycloakErr != nil {
		slog.Warn("failed to initialize Keycloak client, continuing without it", "error", keycloakErr)
		keycloakClient = nil
	}
	tracer, err := jaeger.InitTracer("Booking")
	if err != nil {
		log.Fatal("init tracer", err)
	}
	return &Handler{db, conn, keycloakClient, tracer}, nil
}

func (handler *Handler) GetTracer() trace.Tracer {
	return handler.tracer
}
