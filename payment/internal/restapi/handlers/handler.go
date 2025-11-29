package handlers

import (
	"log"

	"github.com/h4x4d/parking_net/payment/internal/database_service"
	"github.com/h4x4d/parking_net/pkg/client"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	Database *database_service.DatabaseService
	KeyCloak *client.Client
	tracer   trace.Tracer
}

func NewHandler(connStr string) (*Handler, error) {
	db, err := database_service.NewDatabaseService(connStr)
	if err != nil {
		return nil, err
	}
	keycloakClient, keycloakErr := client.NewClient()
	if keycloakErr != nil {
		log.Printf("Warning: failed to initialize Keycloak client, continuing without it: %v", keycloakErr)
		keycloakClient = nil
	}
	tracer, err := jaeger.InitTracer("Payment")
	if err != nil {
		log.Fatal("init tracer", err)
	}
	return &Handler{db, keycloakClient, tracer}, nil
}

func (handler *Handler) GetTracer() trace.Tracer {
	return handler.tracer
}
