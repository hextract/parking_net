package database_service

import (
	"context"
	"fmt"
	"github.com/h4x4d/parking_net/booking/internal/models"
)

func (ds *DatabaseService) GetAll(parkingPlaceID *int64, userID *string) ([]*models.Booking, error) {
	query := "SELECT id FROM bookings WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if parkingPlaceID != nil {
		query += fmt.Sprintf(" AND parking_place_id=$%d", argIndex)
		args = append(args, *parkingPlaceID)
		argIndex++
	}

	if userID != nil {
		query += fmt.Sprintf(" AND user_id=$%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	bookingIdRow, errGetId := ds.pool.Query(context.Background(), query, args...)
	if errGetId != nil {
		return nil, errGetId
	}
	defer bookingIdRow.Close()

	bookings := make([]*models.Booking, 0)

	for bookingIdRow.Next() {
		var bookingId int64
		errScanId := bookingIdRow.Scan(&bookingId)
		if errScanId != nil {
			return nil, errScanId
		}
		booking, errGetBooking := ds.GetByID(bookingId)
		if errGetBooking != nil {
			return nil, errGetBooking
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}
