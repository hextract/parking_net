package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/h4x4d/parking_net/auth/internal/restapi/handlers"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
)

//go:generate swagger generate server --target ../../internal --name ParkingsAuth --spec ../../api/swagger/auth.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.ParkingsAuthAPI) {
}

var authHandler *handlers.Handler

func configureAPI(api *operations.ParkingsAuthAPI) http.Handler {
	var err error
	authHandler, err = handlers.NewHandler()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize auth handler: %v", err))
	}

	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetMetricsHandler == nil {
		api.GetMetricsHandler = operations.GetMetricsHandlerFunc(func(params operations.GetMetricsParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetMetrics has not yet been implemented")
		})
	}
	api.PostChangePasswordHandler = operations.PostChangePasswordHandlerFunc(authHandler.ChangePasswordHandler)
	api.PostLoginHandler = operations.PostLoginHandlerFunc(authHandler.LoginHandler)
	api.PostRegisterHandler = operations.PostRegisterHandlerFunc(authHandler.RegisterHandler)

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
	return handler
}
