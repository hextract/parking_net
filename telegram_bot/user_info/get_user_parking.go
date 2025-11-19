package user_info

import (
	"telegram_bot/models"
)

func (us *UserInfo) GetUserParking(telegramId int64) *models.ParkingPlace {
	parkingPlace, exists := us.Parking[telegramId]
	if !exists {
		parkingPlace = models.NewParkingPlace()
		us.Parking[telegramId] = parkingPlace
	}
	return parkingPlace
}
