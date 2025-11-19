package database_service

import (
	"context"
	"fmt"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"strings"
)

func (ds *DatabaseService) Update(id int64, parkingPlace *models.ParkingPlace) (*models.ParkingPlace, error) {
	query := `UPDATE parking_places SET`
	var fieldNames []string
	var values []interface{}

	if parkingPlace.Address != nil {
		fieldNames = append(fieldNames, fmt.Sprintf("address = $%d", len(values)+1))
		values = append(values, *parkingPlace.Address)
	}
	if parkingPlace.City != nil {
		fieldNames = append(fieldNames, fmt.Sprintf("city = $%d", len(values)+1))
		values = append(values, *parkingPlace.City)
	}
	if parkingPlace.Name != nil {
		fieldNames = append(fieldNames, fmt.Sprintf("name = $%d", len(values)+1))
		values = append(values, *parkingPlace.Name)
	}
	if parkingPlace.HourlyRate != 0 {
		fieldNames = append(fieldNames, fmt.Sprintf("hourly_rate = $%d", len(values)+1))
		values = append(values, parkingPlace.HourlyRate)
	}
	if parkingPlace.ID != 0 {
		fieldNames = append(fieldNames, fmt.Sprintf("id = $%d", len(values)+1))
		values = append(values, parkingPlace.ID)
	}
	if parkingPlace.ParkingType != "" {
		fieldNames = append(fieldNames, fmt.Sprintf("parking_type = $%d", len(values)+1))
		values = append(values, parkingPlace.ParkingType)
	}
	if parkingPlace.Capacity != 0 {
		fieldNames = append(fieldNames, fmt.Sprintf("capacity = $%d", len(values)+1))
		values = append(values, parkingPlace.Capacity)
	}
	query += fmt.Sprintf(" %s WHERE %s RETURNING *", strings.Join(fieldNames, ", "),
		fmt.Sprintf("id = $%d", len(values)+1))
	values = append(values, id)
	fmt.Println(query, values)
	err := ds.pool.QueryRow(context.Background(), query, values...).Scan(&parkingPlace.ID, parkingPlace.Name,
		parkingPlace.City, parkingPlace.Address, &parkingPlace.ParkingType, &parkingPlace.HourlyRate, 
		&parkingPlace.Capacity, &parkingPlace.OwnerID)
	return parkingPlace, err
}
