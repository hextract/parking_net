package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/h4x4d/parking_net/parking/internal/database_service"
	"github.com/h4x4d/parking_net/parking/internal/grpc/gen"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	Database *database_service.DatabaseService
	gen.UnimplementedParkingServer
}

func NewGRPCServer() (*GRPCServer, error) {
	db, err := database_service.NewDatabaseService(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), "db", os.Getenv("POSTGRES_PORT"), os.Getenv("PARKING_DB_NAME")))
	if err != nil {
		return nil, err
	}
	return &GRPCServer{Database: db}, nil
}

func Register(gRPCServer *grpc.Server) {
	server, err := NewGRPCServer()
	if err != nil {
		os.Exit(1)
	}
	gen.RegisterParkingServer(gRPCServer, server)
}

func (serverApi *GRPCServer) GetParkingPlace(
	ctx context.Context, in *gen.ParkingPlaceRequest) (*gen.ParkingPlaceResponse, error) {

	if err := serverApi.validateInternalRequest(ctx); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed")
	}

	if in.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parking place ID")
	}

	tracer := otel.Tracer("Parking")
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md.Get("x-trace-id")) > 0 {
		traceIdString := md.Get("x-trace-id")[0]
		traceId, err := trace.TraceIDFromHex(traceIdString)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid trace ID")
		}
		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceId,
		})
		ctx = trace.ContextWithSpanContext(ctx, spanContext)
	} else {
		ctx = context.Background()
	}
	_, span := tracer.Start(ctx, "get parking place")
	defer span.End()

	parkingPlace, err := serverApi.Database.GetById(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get parking place")
	}
	if parkingPlace == nil {
		return nil, status.Errorf(codes.NotFound, "parking place not found")
	}

	name := ""
	if parkingPlace.Name != nil {
		name = *parkingPlace.Name
	}
	city := ""
	if parkingPlace.City != nil {
		city = *parkingPlace.City
	}
	address := ""
	if parkingPlace.Address != nil {
		address = *parkingPlace.Address
	}

	return &gen.ParkingPlaceResponse{
		Id:          parkingPlace.ID,
		Name:        name,
		City:        city,
		Address:     address,
		ParkingType: parkingPlace.ParkingType,
		HourlyRate:  parkingPlace.HourlyRate,
		Capacity:    parkingPlace.Capacity,
		OwnerId:     parkingPlace.OwnerID,
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
