package api_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"telegram_bot/models"
)

func (s *Service) CreateParkingPlace(parkingPlace *models.ParkingPlace, user *models.User) (bool, error) {
	if parkingPlace == nil {
		return false, errors.New("parking place is nil")
	}

	parkingJSON, errorEncode := json.Marshal(parkingPlace)
	if errorEncode != nil {
		return false, errorEncode
	}

	request, errRequest := s.CreateRequest("POST", s.parkingUrl+"parking/", user)
	if errRequest != nil {
		return false, errRequest
	}
	request.Body = io.NopCloser(bytes.NewBuffer(parkingJSON))

	responseCreate, errCreate := s.client.Do(request)
	if errCreate != nil {
		return false, errCreate
	}
	defer responseCreate.Body.Close()

	if responseCreate.StatusCode != http.StatusOK {
		return false, nil
	}
	return true, nil
}
