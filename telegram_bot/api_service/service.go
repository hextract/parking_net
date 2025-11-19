package api_service

import (
	"net/http"
	"os"
	"telegram_bot/database_service"

	"github.com/h4x4d/parking_net/pkg/client"
)

type Service struct {
	parkingUrl string
	bookingUrl string
	authUrl    string
	database_service.DatabaseService
	client.Client
	client http.Client
}

func NewService() (*Service, error) {
	service := new(Service)
	service.bookingUrl = "http://" + "booking" + ":" + os.Getenv("BOOKING_REST_PORT") + "/"
	service.parkingUrl = "http://" + "parking" + ":" + os.Getenv("PARKING_REST_PORT") + "/"
	service.authUrl = "http://" + "auth" + ":" + os.Getenv("AUTH_REST_PORT") + "/"

	database_pointer, errDatabase := database_service.NewDatabaseService()
	if errDatabase != nil {
		return nil, errDatabase
	}
	service.DatabaseService = *database_pointer
	service.client = http.Client{}

	tokenClient, errorClient := client.NewClient()
	if errorClient != nil {
		return nil, errorClient
	}
	service.Client = *tokenClient

	return service, nil
}
