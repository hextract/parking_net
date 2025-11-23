package restapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	swaggererrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/h4x4d/parking_net/parking/internal/di"
	"github.com/h4x4d/parking_net/parking/internal/models"
	"github.com/h4x4d/parking_net/parking/internal/restapi/handlers"
	"github.com/h4x4d/parking_net/parking/internal/restapi/operations"
	"github.com/h4x4d/parking_net/parking/internal/restapi/operations/instruments"
	"github.com/h4x4d/parking_net/parking/internal/restapi/operations/parking"
	"github.com/h4x4d/parking_net/pkg/client"
	"github.com/h4x4d/parking_net/pkg/middlewares"
)

//go:generate swagger generate server --target ../../internal --name ParkingsParking --spec ../../api/swagger/parking.yaml --principal models.User --exclude-main

func configureFlags(api *operations.ParkingsParkingAPI) {
}

var container *di.Container
var keycloakClient *client.Client
var prometheusMetrics *middlewares.PrometheusMetrics

func configureAPI(api *operations.ParkingsParkingAPI) http.Handler {
	var err error
	container, err = di.NewContainer()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize dependencies: %v", err))
	}

	keycloakClient, err = client.NewClient()
	if err != nil {
		slog.Error("failed to initialize Keycloak client, using mock auth", "error", err)
		keycloakClient = nil
	} else {
		slog.Info("Keycloak client initialized successfully")
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

	api.ParkingCreateParkingHandler = parking.CreateParkingHandlerFunc(container.ParkingHandler.CreateParking)
	api.ParkingGetParkingByIDHandler = parking.GetParkingByIDHandlerFunc(container.ParkingHandler.GetParkingByID)
	api.ParkingGetParkingsHandler = parking.GetParkingsHandlerFunc(container.ParkingHandler.GetParkings)
	api.ParkingUpdateParkingHandler = parking.UpdateParkingHandlerFunc(container.ParkingHandler.UpdateParking)
	api.ParkingDeleteParkingHandler = parking.DeleteParkingHandlerFunc(container.ParkingHandler.DeleteParking)

	api.PreServerShutdown = func() {}
	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

func configureTLS(tlsConfig *tls.Config) {
}

func configureServer(s *http.Server, scheme, addr string) {
}

func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return prometheusMetrics.ApplyMetrics(handler)
}
