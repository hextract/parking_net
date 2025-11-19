package handlers

import (
	"context"
	"fmt"
	"os"

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

	tracer := otel.Tracer("Parking")
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md.Get("x-trace-id")) > 0 {
		traceIdString := md.Get("x-trace-id")[0]
		traceId, err := trace.TraceIDFromHex(traceIdString)
		if err != nil {
			return nil, err
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
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	if parkingPlace == nil {
		return nil, status.Errorf(codes.NotFound, "parking place %d not found", in.Id)
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
