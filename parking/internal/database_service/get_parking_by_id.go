package database_service

import (
	"context"
	"github.com/h4x4d/parking_net/parking/internal/models"
)

func (ds *DatabaseService) GetById(parkingPlaceID int64) (*models.ParkingPlace, error) {
	parkingRow, errGet := ds.pool.Query(context.Background(),
		"SELECT * FROM parking_places WHERE id = $1", parkingPlaceID)
	if errGet != nil {
		return nil, errGet
	}
	if !parkingRow.Next() {
		return nil, nil
	}

	parkingPlace := new(models.ParkingPlace)
	parkingPlace.Name = new(string)
	parkingPlace.City = new(string)
	parkingPlace.Address = new(string)

	err := parkingRow.Scan(&parkingPlace.ID, parkingPlace.Name, parkingPlace.City,
		parkingPlace.Address, &parkingPlace.ParkingType, &parkingPlace.HourlyRate, 
		&parkingPlace.Capacity, &parkingPlace.OwnerID)
	parkingRow.Close()

	return parkingPlace, err
}
