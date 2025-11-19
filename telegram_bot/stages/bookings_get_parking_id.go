package stages

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram_bot/user_info"
)

type BookingsGetParkingIdStage struct {
	InputStages
}

func NewBookingGetParkingIdStage() *BookingsGetParkingIdStage {
	bookingsGetParkingIdStage := new(BookingsGetParkingIdStage)

	bookingsGetParkingIdStage.InputStages = *NewInputStages()

	idInput := InputStage{
		Message: "Введите ID парковки",
		Input: func(userInfo *user_info.UserInfo, telegramId int64, parkingId string) error {
			booking := userInfo.GetUserBooking(telegramId)
			var errId error
			*booking.ParkingPlaceID, errId = strconv.ParseInt(parkingId, 10, 64)
			if errId != nil {
				return errId
			}
			return nil
		},
		Keyboard: tgbotapi.NewRemoveKeyboard(true),
	}

	bookingsGetParkingIdStage.AddStages(idInput)
	return bookingsGetParkingIdStage
}
