package database_service

import (
	"context"
	"fmt"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"go.opentelemetry.io/otel"
	"strings"
	"time"
)

func (ds *DatabaseService) CreateBooking(booking *models.Booking) (*int64, error) {
	query := `INSERT INTO bookings`
	// maybe fieldNames can be placed in common place cause other methods also need this info
	var fieldNames []string
	var fields []string
	var values []interface{}

	if booking.DateFrom != nil {
		fieldNames = append(fieldNames, "date_from")
		date, err := time.Parse("02-01-2006", *booking.DateFrom)
		if err != nil {
			return nil, err
		}
		values = append(values, date.Format(time.DateOnly))
	}
	if booking.DateTo != nil {
		fieldNames = append(fieldNames, "date_to")
		date, err := time.Parse("02-01-2006", *booking.DateTo)
		if err != nil {
			return nil, err
		}
		values = append(values, date.Format(time.DateOnly))
	}
	if booking.ParkingPlaceID != nil {
		fieldNames = append(fieldNames, "parking_place_id")
		values = append(values, booking.ParkingPlaceID)
	}

	if booking.BookingID != 0 {
		fieldNames = append(fieldNames, "booking_id")
		values = append(values, booking.BookingID)
	}

	fieldNames = append(fieldNames, "status")
	values = append(values, booking.Status)
	fieldNames = append(fieldNames, "full_cost")
	values = append(values, booking.FullCost)
	fieldNames = append(fieldNames, "user_id")
	values = append(values, booking.UserID)

	for ind := 0; ind < len(fieldNames); ind++ {
		fields = append(fields, fmt.Sprintf("$%d", ind+1))
	}
	query += fmt.Sprintf(" (%s) VALUES (%s) RETURNING id", strings.Join(fieldNames, ", "),
		strings.Join(fields, ", "))
	errInsert := ds.pool.QueryRow(context.Background(), query, values...).Scan(&booking.BookingID)
	if errInsert != nil {
		return nil, errInsert
	}

	return &booking.BookingID, errInsert
}

func (ds *DatabaseService) Create(ctx context.Context, dateFrom *string, dateTo *string, parkingPlaceID *int64, userID string) (*int64, error) {
	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "create booking in database")
	defer span.End()

	parkingPlace, err := client.GetParkingPlaceById(childCtx, parkingPlaceID)

	if err != nil {
		return nil, err
	}
	dFrom, dateErr1 := time.Parse("02-01-2006", *dateFrom)
	dTo, dateErr2 := time.Parse("02-01-2006", *dateTo)
	if dateErr1 != nil {
		return nil, dateErr1
	}
	if dateErr2 != nil {
		return nil, dateErr2
	}
	hours := int64(dTo.Sub(dFrom).Hours())
	cost := parkingPlace.HourlyRate * hours

	booking := &models.Booking{
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		ParkingPlaceID:  parkingPlaceID,
		FullCost:        cost,
		Status:          "Waiting",
		UserID:          userID,
	}

	return ds.CreateBooking(booking)
}
