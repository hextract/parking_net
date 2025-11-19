package api_service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"telegram_bot/models"
)

func (s *Service) GetParkingPlaces(user *models.User) ([]models.ParkingPlace, error) {
	request, errRequest := s.CreateRequest("GET", s.parkingUrl+"parking/", user)
	if errRequest != nil {
		return nil, errRequest
	}

	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	parkingJSON, errJSON := ioutil.ReadAll(response.Body)
	if errJSON != nil {
		return nil, errJSON
	}

	var parkingPlaces []models.ParkingPlace

	errDecode := json.Unmarshal(parkingJSON, &parkingPlaces)
	if errDecode != nil {
		return nil, errDecode
	}
	return parkingPlaces, nil
}
