package database_service

import (
	"context"
	"github.com/go-openapi/strfmt"
	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func (ds *DatabaseService) GetByID(BookingID int64) (*models.Booking, error) {
	bookingRow, errGet := ds.pool.Query(context.Background(),
		"SELECT * FROM bookings WHERE id = $1", BookingID)
	if errGet != nil {
		return nil, errGet
	}
	defer bookingRow.Close()

	if !bookingRow.Next() {
		return nil, nil
	}

	booking := new(models.Booking)
	booking.ParkingPlaceID = new(int64)

	from := new(pgtype.Timestamp)
	to := new(pgtype.Timestamp)

	errBooking := bookingRow.Scan(&booking.BookingID, from,
		to, booking.ParkingPlaceID, &booking.FullCost, &booking.Status, &booking.UserID)

	fromDT := strfmt.DateTime(from.Time)
	toDT := strfmt.DateTime(to.Time)
	booking.DateFrom = &fromDT
	booking.DateTo = &toDT
	return booking, errBooking
}
