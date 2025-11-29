package database_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
)

func (ds *DatabaseService) Update(ctx context.Context, bookingId int64, booking *models.Booking) (*models.Booking, error) {
	query := `UPDATE bookings SET`
	var settings []string
	var values []interface{}

	tracer := otel.Tracer("Booking")
	ctx, span := tracer.Start(ctx, "update")
	defer span.End()

	if booking.DateFrom != nil {
		settings = append(settings, fmt.Sprintf("date_from = $%d", len(values)+1))
		values = append(values, time.Time(*booking.DateFrom))
	}

	if booking.DateTo != nil {
		settings = append(settings, fmt.Sprintf("date_to = $%d", len(values)+1))
		values = append(values, time.Time(*booking.DateTo))
	}

	if booking.ParkingPlaceID != nil {
		settings = append(settings, fmt.Sprintf("parking_place_id = $%d", len(values)+1))
		values = append(values, *booking.ParkingPlaceID)
	}
	if booking.FullCost == 0 {
		if booking.ParkingPlaceID == nil {
			return nil, fmt.Errorf("parking place ID is required")
		}
		if err := utils.ValidateParkingPlaceID(booking.ParkingPlaceID); err != nil {
			return nil, fmt.Errorf("invalid parking place ID")
		}

		parkingPlace, err := client.GetParkingPlaceById(ctx, booking.ParkingPlaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parking place")
		}

		if booking.DateFrom == nil || booking.DateTo == nil {
			return nil, fmt.Errorf("dates cannot be nil")
		}

		dFrom := time.Time(*booking.DateFrom)
		dTo := time.Time(*booking.DateTo)
		if err := utils.ValidateDateRange(&dFrom, &dTo); err != nil {
			return nil, err
		}

		hours := dTo.Sub(dFrom).Hours()
		booking.FullCost = int64(float64(parkingPlace.HourlyRate) * hours)
		if err := utils.ValidateFullCost(booking.FullCost); err != nil {
			return nil, fmt.Errorf("calculated cost exceeds maximum")
		}
	} else {
		if err := utils.ValidateFullCost(booking.FullCost); err != nil {
			return nil, fmt.Errorf("invalid full cost")
		}
	}

	settings = append(settings, fmt.Sprintf("full_cost = $%d", len(values)+1))
	values = append(values, booking.FullCost)

	if booking.Status != "" {
		settings = append(settings, fmt.Sprintf("status = $%d", len(values)+1))
		values = append(values, booking.Status)
	}

	if booking.UserID != "" {
		settings = append(settings, fmt.Sprintf("user_id = $%d", len(values)+1))
		values = append(values, booking.UserID)
	}

	query += fmt.Sprintf(" %s WHERE %s RETURNING *", strings.Join(settings, ", "),
		fmt.Sprintf("id = $%d", len(values)+1))
	values = append(values, bookingId)

	from := new(pgtype.Timestamp)
	to := new(pgtype.Timestamp)

	errUpdate := ds.pool.QueryRow(context.Background(), query, values...).Scan(&booking.BookingID, from,
		to, booking.ParkingPlaceID, &booking.FullCost, &booking.Status, &booking.UserID)

	fromDT := strfmt.DateTime(from.Time)
	toDT := strfmt.DateTime(to.Time)
	booking.DateFrom = &fromDT
	booking.DateTo = &toDT
	return booking, errUpdate
}
