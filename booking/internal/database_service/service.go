package database_service

import (
	"context"
	"fmt"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	pool        *pgxpool.Pool
	parkingPool *pgxpool.Pool
}

func NewDatabaseService(connStr string) (*DatabaseService, error) {
	result := new(DatabaseService)
	newPool, errPool := pgxpool.New(context.Background(), connStr)
	if errPool != nil {
		return nil, errPool
	}
	result.pool = newPool

	parkingConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"db",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("PARKING_DB_NAME"),
	)
	parkingPool, errParkingPool := pgxpool.New(context.Background(), parkingConnStr)
	if errParkingPool != nil {
		return nil, fmt.Errorf("failed to create parking pool: %w", errParkingPool)
	}
	result.parkingPool = parkingPool

	return result, nil
}
