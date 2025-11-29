package api_service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

type BalanceResponse struct {
	Balance  *int64  `json:"balance"`
	Currency *string `json:"currency"`
	UserID   *string `json:"user_id"`
}

func (s *Service) GetBalance(token string) (*BalanceResponse, error) {
	paymentUrl := "http://payment:" + os.Getenv("PAYMENT_REST_PORT") + "/payment/balance"
	req, err := http.NewRequest("GET", paymentUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("api_key", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get balance")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var balance BalanceResponse
	if err := json.Unmarshal(body, &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}
