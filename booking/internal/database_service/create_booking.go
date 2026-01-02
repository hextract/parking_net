package database_service

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/booking/internal/grpc/client"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/utils"
	"github.com/jackc/pgx/v5"
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

func (ds *DatabaseService) Create(ctx context.Context, dateFrom *strfmt.DateTime, dateTo *strfmt.DateTime, parkingPlaceID *int64, userID string, services []BookingService) (*int64, error) {
	if err := utils.ValidateUserID(userID); err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	if err := utils.ValidateParkingPlaceID(parkingPlaceID); err != nil {
		return nil, fmt.Errorf("invalid parking place ID")
	}

	if dateFrom == nil || dateTo == nil {
		return nil, fmt.Errorf("dates cannot be nil")
	}

	dFrom := time.Time(*dateFrom)
	dTo := time.Time(*dateTo)
	if err := utils.ValidateDateRange(&dFrom, &dTo); err != nil {
		return nil, err
	}

	tracer := otel.Tracer("Booking")
	childCtx, span := tracer.Start(ctx, "create booking in database")
	defer span.End()

	parkingPlace, err := client.GetParkingPlaceById(childCtx, parkingPlaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parking place")
	}

	hours := dTo.Sub(dFrom).Hours()
	cost := int64(float64(parkingPlace.HourlyRate) * hours)

	servicesCost := int64(0)
	if len(services) > 0 && ds.parkingPool != nil {
		for i := range services {
			query := `SELECT id, parking_id, name, description, price
				FROM services WHERE id = $1`
			
			var serviceInfo struct {
				ID          int64
				ParkingID   int64
				Name        string
				Description string
				Price       int64
			}
			
			err := ds.parkingPool.QueryRow(ctx, query, services[i].ServiceID).Scan(
				&serviceInfo.ID,
				&serviceInfo.ParkingID,
				&serviceInfo.Name,
				&serviceInfo.Description,
				&serviceInfo.Price,
			)
			
			if err != nil {
				if err == pgx.ErrNoRows {
					return nil, fmt.Errorf("service %d not found", services[i].ServiceID)
				}
				return nil, fmt.Errorf("failed to get service %d: %w", services[i].ServiceID, err)
			}
			
			if serviceInfo.ParkingID != *parkingPlaceID {
				return nil, fmt.Errorf("service %d does not belong to parking place %d", services[i].ServiceID, *parkingPlaceID)
			}
			
			services[i].Price = serviceInfo.Price
			servicesCost += serviceInfo.Price * services[i].Quantity
		}
	}

	totalCost := cost + servicesCost
	if err := utils.ValidateFullCost(totalCost); err != nil {
		return nil, fmt.Errorf("calculated cost exceeds maximum")
	}

	booking := &models.Booking{
		DateFrom:        dateFrom,
		DateTo:          dateTo,
		ParkingPlaceID:  parkingPlaceID,
		FullCost:        totalCost,
		Status:          "Waiting",
		UserID:          userID,
	}

	bookingID, err := ds.CreateBooking(booking)
	if err != nil {
		return nil, err
	}

	if len(services) > 0 {
		err = ds.AddBookingServices(ctx, *bookingID, services)
		if err != nil {
			return nil, fmt.Errorf("failed to add booking services: %w", err)
		}
	}

	return bookingID, nil
}
