package data_representation

import (
	"strconv"
	"telegram_bot/models"
)

func GetParkingPlace(parkingPlace *models.ParkingPlace) string {
	var result string
	result += "Parking place information;\n"
	result += "Parking ID: " + "\"" + strconv.FormatInt(parkingPlace.ID, 10) + "\"" + ";\n"
	result += "Name: " + "\"" + *parkingPlace.Name + "\"" + "\n"
	result += "Address: " + "\"" + *parkingPlace.City + ", " + *parkingPlace.Address + "\"" + ";\n"
	result += "Hourly rate: " + strconv.FormatInt(parkingPlace.HourlyRate, 10) + " USD" + ";\n"
	result += "Capacity: " + strconv.FormatInt(parkingPlace.Capacity, 10) + " spaces" + ";\n"
	result += "Parking type: " + "\"" + parkingPlace.ParkingType + "\"" + ";\n"
	result += "Created by owner with ID " + "\"" + parkingPlace.OwnerID[:12] + "..." + "\"" + ";\n"

	return result
}
