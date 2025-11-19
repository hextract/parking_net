package database_service

import (
	"context"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"go.opentelemetry.io/otel"
)

func (ds *DatabaseService) CheckOwnership(ctx context.Context, BookingID int64, user *models.User) (bool, error) {
	if user == nil {
		return false, nil
	}
	if user.Role == "driver" {
		booking, err := ds.GetByID(BookingID)
		if err != nil {
			return false, err
		}
		if booking == nil {
			return false, nil
		}
		return booking.UserID == user.UserID, nil
	}
	booking, err := ds.GetByID(BookingID)
	if err != nil {
		return false, err
	}
	if booking == nil {
		return false, nil
	}

	tracer := otel.Tracer("Booking")
	ctx, span := tracer.Start(ctx, "check ownership db")
	defer span.End()

	parkingPlace, parkingErr := client.GetParkingPlaceById(ctx, booking.ParkingPlaceID)
	if parkingErr != nil {
		return false, parkingErr
	}
	return parkingPlace.OwnerID == user.UserID, nil
}
