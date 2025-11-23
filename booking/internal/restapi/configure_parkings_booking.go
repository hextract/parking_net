package restapi

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	swaggererrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/h4x4d/parking_net/booking/internal/models"
	"github.com/h4x4d/parking_net/booking/internal/restapi/handlers"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/driver"
	"github.com/h4x4d/parking_net/booking/internal/restapi/operations/instruments"
	"github.com/h4x4d/parking_net/pkg/client"
	"github.com/h4x4d/parking_net/pkg/middlewares"
)

//go:generate swagger generate server --target ../../internal --name ParkingsBooking --spec ../../api/swagger/booking.yaml --principal models.User --exclude-main

func configureFlags(api *operations.ParkingsBookingAPI) {
}

var keycloakClient *client.Client
var bookingHandler *handlers.Handler
var prometheusMetrics *middlewares.PrometheusMetrics

func configureAPI(api *operations.ParkingsBookingAPI) http.Handler {
	var err error
	keycloakClient, err = client.NewClient()
	if err != nil {
		slog.Warn("failed to initialize Keycloak client, using mock auth", "error", err)
		keycloakClient = nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"db",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("BOOKING_DB_NAME"),
	)
	bookingHandler, err = handlers.NewHandler(connStr)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize booking handler: %v", err))
	}

	prometheusMetrics = middlewares.NewPrometheusMetrics()

	api.ServeError = swaggererrors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.APIKeyAuth = func(token string) (*models.User, error) {
		if token == "" {
			return nil, errors.New("token required")
		}

		if keycloakClient == nil {
			slog.Warn("Keycloak client not available, using mock auth")
			return &models.User{
				UserID:     "test-user-123",
				Role:       "owner",
				TelegramID: 123456789,
			}, nil
		}

		keycloakUser, err := keycloakClient.CheckToken(context.Background(), token)
		if err != nil {
			slog.Error("failed to validate token", "error", err)
			return nil, fmt.Errorf("invalid token: %w", err)
		}

		return &models.User{
			UserID:     keycloakUser.UserID,
			Role:       keycloakUser.Role,
			TelegramID: keycloakUser.TelegramID,
		}, nil
	}

	api.InstrumentsGetMetricsHandler = instruments.GetMetricsHandlerFunc(handlers.MetricsHandler)
	api.DriverCreateBookingHandler = driver.CreateBookingHandlerFunc(bookingHandler.CreateBooking)
	api.DriverGetBookingHandler = driver.GetBookingHandlerFunc(bookingHandler.GetBooking)
	api.DriverGetBookingByIDHandler = driver.GetBookingByIDHandlerFunc(bookingHandler.GetBookingByID)
	api.DriverUpdateBookingHandler = driver.UpdateBookingHandlerFunc(bookingHandler.UpdateBooking)
	api.DriverDeleteBookingHandler = driver.DeleteBookingHandlerFunc(bookingHandler.DeleteBooking)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func configureTLS(tlsConfig interface{}) {
}

func configureServer(s *http.Server, scheme, addr string) {
}

func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return prometheusMetrics.ApplyMetrics(handler)
}
