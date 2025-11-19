package stages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram_bot/api_service"
	"telegram_bot/user_info"
)

type CreateParkingStage struct {
	InputStages
}

func (cps *CreateParkingStage) Finish(userInfo *user_info.UserInfo, telegramId int64, apiService *api_service.Service) (bool, error) {
	user := userInfo.GetUserData(telegramId)
	parkingPlace := userInfo.GetUserParking(telegramId)
	parkingPlace.OwnerID = user.UserID
	parkingPlace.ID = 0
	return apiService.CreateParkingPlace(parkingPlace, user)
}

func NewCreateParkingStage() *CreateParkingStage {
	createParkingStage := new(CreateParkingStage)

	createParkingStage.InputStages = *NewInputStages()

	nameInput := InputStage{
		Message: "Введите название",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, name string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			*parkingPlace.Name = name
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	addressInput := InputStage{
		Message: "Введите адрес",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, address string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			*parkingPlace.Address = address
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	cityInput := InputStage{
		Message: "Введите город",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, city string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			*parkingPlace.City = city
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	hourlyRateInput := InputStage{
		Message: "Введите стоимость за час (в рублях)",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, rate string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			var errorParse error
			parkingPlace.HourlyRate, errorParse = strconv.ParseInt(rate, 10, 64)
			if errorParse != nil {
				return errorParse
			}
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	capacityInput := InputStage{
		Message: "Введите вместимость (количество мест)",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, capacity string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			var errorParse error
			parkingPlace.Capacity, errorParse = strconv.ParseInt(capacity, 10, 64)
			if errorParse != nil {
				return errorParse
			}
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	parkingTypeInput := InputStage{
		Message: "Выберите тип парковки",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, parkingType string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			parkingPlace.ParkingType = parkingType
			return nil
		},
		Keyboard: tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("outdoor"),
				tgbotapi.NewKeyboardButton("covered"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("underground"),
				tgbotapi.NewKeyboardButton("multi-level"),
			)),
	}

	createParkingStage.AddStages(nameInput, addressInput, cityInput, hourlyRateInput, capacityInput, parkingTypeInput)
	return createParkingStage
}
