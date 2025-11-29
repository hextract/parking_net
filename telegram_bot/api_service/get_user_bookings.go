package api_service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"telegram_bot/models"
)

func (s *Service) GetUserBookings(user *models.User) ([]models.Booking, error) {
	if user == nil || user.UserID == "" {
		return nil, fmt.Errorf("invalid user")
	}

	path := s.bookingUrl + "booking/"

	urlObject, errUrl := url.Parse(path)
	if errUrl != nil {
		return nil, fmt.Errorf("failed to parse URL")
	}

	params := url.Values{}
	params.Add("user_id", user.UserID)

	urlObject.RawQuery = params.Encode()

	request, errRequest := s.CreateRequest("GET", urlObject.String(), user)
	if errRequest != nil {
		return nil, fmt.Errorf("failed to create request")
	}

	responseBookings, errBookings := s.httpClient.Do(request)
	if errBookings != nil {
		return nil, errBookings
	}
	defer responseBookings.Body.Close()

	if responseBookings.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get bookings")
	}

	bookingsJSON, errJSON := ioutil.ReadAll(responseBookings.Body)
	if errJSON != nil {
		return nil, errJSON
	}

	var bookings []models.Booking

	errDecode := json.Unmarshal(bookingsJSON, &bookings)
	if errDecode != nil {
		return nil, errDecode
	}
	return bookings, nil
}
