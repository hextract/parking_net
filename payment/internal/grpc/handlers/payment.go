package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/h4x4d/parking_net/payment/internal/database_service"
	"github.com/h4x4d/parking_net/payment/internal/grpc/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	Database *database_service.DatabaseService
	gen.UnimplementedPaymentServer
	tracer trace.Tracer
}

func NewGRPCServer() (*GRPCServer, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"db",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("PAYMENT_DB_NAME"),
	)
	db, err := database_service.NewDatabaseService(connStr)
	if err != nil {
		return nil, err
	}
	tracer := otel.Tracer("Payment-gRPC")
	return &GRPCServer{Database: db, tracer: tracer}, nil
}

func Register(gRPCServer *grpc.Server) {
	server, err := NewGRPCServer()
	if err != nil {
		os.Exit(1)
	}
	gen.RegisterPaymentServer(gRPCServer, server)
}

func (s *GRPCServer) ProcessTransaction(ctx context.Context, req *gen.TransactionRequest) (*gen.TransactionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "ProcessTransaction")
	defer span.End()

	result, err := s.Database.ProcessTransaction(ctx, req.BookingId, req.DriverId, req.OwnerId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process transaction: %v", err)
	}

	if result.Status == "failed" {
		return &gen.TransactionResponse{
			TransactionId: result.TransactionID,
			Status:        result.Status,
			Message:       result.Message,
		}, nil
	}

	return &gen.TransactionResponse{
		TransactionId: result.TransactionID,
		Status:        result.Status,
		Message:       result.Message,
	}, nil
}

func (s *GRPCServer) ProcessRefund(ctx context.Context, req *gen.RefundRequest) (*gen.TransactionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "ProcessRefund")
	defer span.End()

	result, err := s.Database.ProcessRefund(ctx, req.BookingId, req.DriverId, req.OwnerId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process refund: %v", err)
	}

	if result.Status == "failed" {
		return &gen.TransactionResponse{
			TransactionId: result.TransactionID,
			Status:        result.Status,
			Message:       result.Message,
		}, nil
	}

	return &gen.TransactionResponse{
		TransactionId: result.TransactionID,
		Status:        result.Status,
		Message:       result.Message,
	}, nil
}
