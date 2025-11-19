package api_service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"telegram_bot/models"
)

func (s *Service) GetParkingPlaceByID(parkingID int64, user *models.User) (*models.ParkingPlace, error) {
	path := s.parkingUrl + "parking/" + strconv.Itoa(int(parkingID))
	request, errRequest := s.CreateRequest("GET", path, user)
	if errRequest != nil {
		return nil, errRequest
	}

	response, err := s.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("get parking place by id failed")
	}

	parkingJSON, errJSON := ioutil.ReadAll(response.Body)
	if errJSON != nil {
		return nil, errJSON
	}

	parkingPlace := new(models.ParkingPlace)

	errDecode := json.Unmarshal(parkingJSON, &parkingPlace)
	if errDecode != nil {
		return nil, errDecode
	}
	return parkingPlace, nil
}
