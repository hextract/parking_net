package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/h4x4d/parking_net/auth/internal/restapi/handlers"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/pkg/middlewares"
)

//go:generate swagger generate server --target ../../internal --name ParkingsAuth --spec ../../api/swagger/auth.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.ParkingsAuthAPI) {
}

var authHandler *handlers.Handler
var prometheusMetrics *middlewares.PrometheusMetrics

func configureAPI(api *operations.ParkingsAuthAPI) http.Handler {
	var err error
	authHandler, err = handlers.NewHandler()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize auth handler: %v", err))
	}

	prometheusMetrics = middlewares.NewPrometheusMetrics()

	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetAuthMetricsHandler = operations.GetAuthMetricsHandlerFunc(handlers.MetricsHandler)
	api.GetAuthMeHandler = operations.GetAuthMeHandlerFunc(authHandler.GetMeHandler)
	api.PostAuthChangePasswordHandler = operations.PostAuthChangePasswordHandlerFunc(authHandler.ChangePasswordHandler)
	api.PostAuthLoginHandler = operations.PostAuthLoginHandlerFunc(authHandler.LoginHandler)
	api.PostAuthRegisterHandler = operations.PostAuthRegisterHandlerFunc(authHandler.RegisterHandler)

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
