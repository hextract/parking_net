package stages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram_bot/user_info"
)

type ParkingGetParkingIdStage struct {
	InputStages
}

func NewParkingGetParkingIdStage() *ParkingGetParkingIdStage {
	parkingGetParkingIdStage := new(ParkingGetParkingIdStage)

	parkingGetParkingIdStage.InputStages = *NewInputStages()

	idInput := InputStage{
		Message: "Введите ID парковки",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, parkingId string) error {
			parkingPlace := userInfo.GetUserParking(telegramId)
			var errId error
			parkingPlace.ID, errId = strconv.ParseInt(parkingId, 10, 64)
			if errId != nil {
				return errId
			}
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	parkingGetParkingIdStage.AddStages(idInput)
	return parkingGetParkingIdStage
}
