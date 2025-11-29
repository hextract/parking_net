package api_service

import (
	"telegram_bot/models"
)

func (s *Service) Logout(user *models.User) error {
	telegramID := user.TelegramID
	err := s.DatabaseService.SetToken(telegramID, &models.ApiToken{Value: ""})
	if err != nil {
		return err
	}
	user = new(models.User)
	user.TelegramID = telegramID
	return nil
}
