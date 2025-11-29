package api_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"telegram_bot/models"
)

func (s *Service) GetUserParkings(user *models.User) ([]models.ParkingPlace, error) {
	if user == nil || user.UserID == "" {
		return nil, fmt.Errorf("invalid user")
	}

	path := s.parkingUrl + "parking/"

	urlObject, errUrl := url.Parse(path)
	if errUrl != nil {
		return nil, fmt.Errorf("failed to parse URL")
	}

	params := url.Values{}
	params.Add("owner_id", user.UserID)

	urlObject.RawQuery = params.Encode()

	request, errRequest := s.CreateRequest("GET", urlObject.String(), user)
	if errRequest != nil {
		return nil, fmt.Errorf("failed to create request")
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get parkings")
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
