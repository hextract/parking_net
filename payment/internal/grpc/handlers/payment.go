package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/h4x4d/parking_net/payment/internal/database_service"
	"github.com/h4x4d/parking_net/payment/internal/grpc/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	if err := s.validateInternalRequest(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	ctx, span := s.tracer.Start(ctx, "ProcessTransaction")
	defer span.End()

	result, err := s.Database.ProcessTransaction(ctx, req.BookingId, req.DriverId, req.OwnerId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "transaction processing failed")
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
	if err := s.validateInternalRequest(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	ctx, span := s.tracer.Start(ctx, "ProcessRefund")
	defer span.End()

	result, err := s.Database.ProcessRefund(ctx, req.BookingId, req.DriverId, req.OwnerId, req.Amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refund processing failed")
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

func (s *GRPCServer) validateInternalRequest(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("no metadata provided")
	}

	internalToken := os.Getenv("INTERNAL_SERVICE_TOKEN")
	if internalToken == "" {
		return nil
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return fmt.Errorf("no authorization header")
	}

	authHeader := strings.TrimSpace(authHeaders[0])
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return fmt.Errorf("invalid authorization format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != internalToken {
		return fmt.Errorf("invalid token")
	}

	return nil
}
