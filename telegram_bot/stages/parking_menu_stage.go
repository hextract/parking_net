package stages

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ParkingMenuStage struct {
	keyboard interface{}
	message  string
}

func (ps *ParkingMenuStage) ConfigureMessage(message *tgbotapi.MessageConfig) {
	message.Text = ps.message
	message.ReplyMarkup = ps.keyboard
}

func NewParkingMenuStage(role string) *ParkingMenuStage {
	parkingMenuStage := new(ParkingMenuStage)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Получить все"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Получить по ID"),
		))

	if role == "owner" {
		keyboard.Keyboard = append(keyboard.Keyboard,
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Создать парковку"),
			))
	}

	keyboard.Keyboard = append(keyboard.Keyboard,
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Вернутся обратно"),
		))

	parkingMenuStage.keyboard = keyboard

	parkingMenuStage.message = "Меню парковок"
	return parkingMenuStage
}
