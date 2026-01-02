package di

import (
	"context"
	"fmt"
	"os"

	"github.com/h4x4d/parking_net/parking/internal/handlers"
	"github.com/h4x4d/parking_net/parking/internal/repository"
	"github.com/h4x4d/parking_net/parking/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	ParkingHandler *handlers.ParkingHandler
	ServiceHandler *handlers.ServiceHandler
}

func NewContainer() (*Container, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"db",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("PARKING_DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	parkingRepo := repository.NewPostgresParkingRepository(pool)
	parkingSvc := service.NewParkingService(parkingRepo)

	parkingHandler, err := handlers.NewParkingHandler(parkingSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to create parking handler: %w", err)
	}

	serviceRepo := repository.NewPostgresServiceRepository(pool)
	serviceSvc := service.NewServiceService(serviceRepo, parkingRepo)

	serviceHandler, err := handlers.NewServiceHandler(serviceSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to create service handler: %w", err)
	}

	return &Container{
		ParkingHandler: parkingHandler,
		ServiceHandler: serviceHandler,
	}, nil
}
