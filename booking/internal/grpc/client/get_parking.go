package client

import (
	"context"
	"os"

	"github.com/h4x4d/parking_net/booking/internal/grpc/gen"
	"github.com/h4x4d/parking_net/booking/internal/grpc/utils"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/metadata"
)

func GetParkingPlaceById(ctx context.Context, parkingPlaceId *int64) (*models.ParkingPlace, error) {
	conn, err := utils.ConnectToParking()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "booking request get parking place")
	defer span.End()

	internalToken := os.Getenv("INTERNAL_SERVICE_TOKEN")
	if internalToken != "" {
		childCtx = metadata.AppendToOutgoingContext(childCtx, "authorization", "Bearer "+internalToken)
	}

	client := gen.NewParkingClient(conn)

	parkingResp, err := client.GetParkingPlace(childCtx, &gen.ParkingPlaceRequest{Id: *parkingPlaceId})
	if err != nil {
		return nil, err
	}
	parkingPlace := models.ParkingPlace{
		ID:         parkingResp.Id,
		Name:       &parkingResp.Name,
		City:       &parkingResp.City,
		Address:    &parkingResp.Address,
		HourlyRate: parkingResp.HourlyRate,
		Capacity:   parkingResp.Capacity,
		ParkingType: parkingResp.ParkingType,
		OwnerID:    parkingResp.OwnerId,
	}
	return &parkingPlace, err
}
