package api_service

import (
	"fmt"
	"net/http"
	"telegram_bot/models"
)

func (s *Service) CreateRequest(method string, path string, user *models.User) (*http.Request, error) {
	request, errRequest := http.NewRequest(method, path, nil)
	if errRequest != nil {
		return nil, errRequest
	}

	apiToken, errToken := s.DatabaseService.GetToken(user.TelegramID)
	if errToken != nil {
		return nil, errToken
	}
	if apiToken == nil {
		return nil, fmt.Errorf("token not found")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("api_key", apiToken.Value)
	return request, nil
}
