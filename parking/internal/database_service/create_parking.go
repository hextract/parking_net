package database_service

import (
	"context"
	"fmt"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"strings"
)

func (ds *DatabaseService) Create(parkingPlace *models.ParkingPlace, user *models.User) (*int64, error) {
	query := `INSERT INTO parking_places`
	var fieldNames []string
	var fields []string
	var values []interface{}

	if parkingPlace.Address != nil {
		fieldNames = append(fieldNames, "address")
		values = append(values, parkingPlace.Address)
	}
	if parkingPlace.City != nil {
		fieldNames = append(fieldNames, "city")
		values = append(values, parkingPlace.City)
	}
	if parkingPlace.Name != nil {
		fieldNames = append(fieldNames, "name")
		values = append(values, parkingPlace.Name)
	}
	if parkingPlace.ID != 0 {
		fieldNames = append(fieldNames, "id")
		values = append(values, parkingPlace.ID)
	}
	fieldNames = append(fieldNames, "hourly_rate")
	values = append(values, parkingPlace.HourlyRate)

	if parkingPlace.ParkingType != "" {
		fieldNames = append(fieldNames, "parking_type")
		values = append(values, parkingPlace.ParkingType)
	}

	fieldNames = append(fieldNames, "capacity")
	values = append(values, parkingPlace.Capacity)

	fieldNames = append(fieldNames, "owner_id")
	values = append(values, user.UserID)

	for ind := 0; ind < len(fieldNames); ind++ {
		fields = append(fields, fmt.Sprintf("$%d", ind+1))
	}
	query += fmt.Sprintf(" (%s) VALUES (%s) RETURNING id", strings.Join(fieldNames, ", "),
		strings.Join(fields, ", "))
	errInsert := ds.pool.QueryRow(context.Background(), query, values...).Scan(&parkingPlace.ID)
	if errInsert != nil {
		return nil, errInsert
	}

	return &parkingPlace.ID, errInsert
}
