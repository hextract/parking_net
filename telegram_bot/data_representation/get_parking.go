package data_representation

import (
	"strconv"
	"telegram_bot/models"
)

func GetParkingPlace(parkingPlace *models.ParkingPlace) string {
	var result string
	result += "Информация о парковке;\n"
	result += "ID парковки: " + "\"" + strconv.FormatInt(parkingPlace.ID, 10) + "\"" + ";\n"
	result += "Название: " + "\"" + *parkingPlace.Name + "\"" + "\n"
	result += "Адрес: " + "\"" + *parkingPlace.City + ", " + *parkingPlace.Address + "\"" + ";\n"
	result += "Стоимость за час: " + strconv.FormatInt(parkingPlace.HourlyRate, 10) + "₽" + ";\n"
	result += "Вместимость: " + strconv.FormatInt(parkingPlace.Capacity, 10) + " мест" + ";\n"
	result += "Тип парковки: " + "\"" + parkingPlace.ParkingType + "\"" + ";\n"
	result += "Создана владельцем с ID " + "\"" + parkingPlace.OwnerID[:12] + "..." + "\"" + ";\n"

	return result
}
