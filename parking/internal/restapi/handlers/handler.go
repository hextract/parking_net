package handlers

import (
	"github.com/h4x4d/parking_net/parking/internal/database_service"
	"github.com/h4x4d/parking_net/pkg/jaeger"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type Handler struct {
	Database *database_service.DatabaseService
	tracer   trace.Tracer
}

func NewHandler(connStr string) (*Handler, error) {
	db, err := database_service.NewDatabaseService(connStr)
	if err != nil {
		return nil, err
	}

	tracer, err := jaeger.InitTracer("Parking")
	if err != nil {
		log.Fatal("init tracer", err)
	}

	return &Handler{db, tracer}, nil
}
