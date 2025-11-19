package database_service

import (
	"context"
)

func (ds *DatabaseService) Exists(parkingPlaceID int64) (bool, error) {
	row, errGet := ds.pool.Query(context.Background(),
		"SELECT id FROM parking_places WHERE id = $1", parkingPlaceID)
	if errGet != nil {
		return false, errGet
	}
	status := row.Next()
	row.Close()
	return status, nil
}
