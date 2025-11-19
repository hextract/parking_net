package database_service

import (
	"context"
	"fmt"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"strings"
)

func (ds *DatabaseService) GetAll(city *string, parkingType *string, name *string) ([]*models.ParkingPlace, error) {
	query := `SELECT * FROM parking_places`
	var clauses []string
	var args []interface{}

	if city != nil {
		clauses = append(clauses, fmt.Sprintf("city = $%d", len(clauses)+1))
		args = append(args, *city)
	}
	if parkingType != nil {
		clauses = append(clauses, fmt.Sprintf("parking_type = $%d", len(clauses)+1))
		args = append(args, *parkingType)
	}
	if name != nil {
		clauses = append(clauses, fmt.Sprintf("NAME LIKE $%d", len(clauses)+1))
		args = append(args, "%"+*name+"%")
	}
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	rows, errQuery := ds.pool.Query(context.Background(), query, args...)
	if errQuery != nil {
		return nil, errQuery
	}

	result := make([]*models.ParkingPlace, 0)
	for rows.Next() {
		parkingPlace := new(models.ParkingPlace)
		parkingPlace.Name = new(string)
		parkingPlace.City = new(string)
		parkingPlace.Address = new(string)

		err := rows.Scan(&parkingPlace.ID, parkingPlace.Name, parkingPlace.City,
			parkingPlace.Address, &parkingPlace.ParkingType, &parkingPlace.HourlyRate, 
			&parkingPlace.Capacity, &parkingPlace.OwnerID)
		if err != nil {
			return nil, err
		}
		result = append(result, parkingPlace)
	}
	rows.Close()
	return result, nil
}
