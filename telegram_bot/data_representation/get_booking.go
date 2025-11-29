package data_representation

import (
	"strconv"
	"telegram_bot/models"
)

func GetBooking(booking *models.Booking) string {
	var result string
	result += "Booking information\n"
	result += "Booking ID: " + "\"" + strconv.FormatInt(booking.BookingID, 10) + "\"" + ";\n"
	result += "Booking period: " + "\"" + *booking.DateFrom + " - " + *booking.DateTo + "\"" + ";\n"
	result += "Booking status: " + "\"" + booking.Status + "\"" + ";\n"
	result += "Created by user with ID " + "\"" + booking.UserID[:12] + "..." + "\"" + ";\n"
	result += "Related to parking place with ID " + "\"" + strconv.FormatInt(*booking.ParkingPlaceID, 10) + "\"" + ";\n"

	return result
}
