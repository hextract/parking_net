package database_service

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
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
		values = append(values, time.Time(*booking.DateFrom))
	}
	if booking.DateTo != nil {
		fieldNames = append(fieldNames, "date_to")
		values = append(values, time.Time(*booking.DateTo))
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

func (ds *DatabaseService) Create(ctx context.Context, dateFrom *strfmt.DateTime, dateTo *strfmt.DateTime, parkingPlaceID *int64, userID string) (*int64, error) {
	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "create booking in database")
	defer span.End()

	parkingPlace, err := client.GetParkingPlaceById(childCtx, parkingPlaceID)

	if err != nil {
		return nil, err
	}
	dFrom := time.Time(*dateFrom)
	dTo := time.Time(*dateTo)
	if dFrom.After(dTo) || dFrom.Equal(dTo) {
		return nil, fmt.Errorf("date_from must be before date_to")
	}
	hours := dTo.Sub(dFrom).Hours()
	cost := int64(float64(parkingPlace.HourlyRate) * hours)

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
