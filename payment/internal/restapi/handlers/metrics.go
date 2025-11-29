package handlers

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/payment/internal/restapi/operations/instruments"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func MetricsHandler(params instruments.GetMetricsParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, _ runtime.Producer) {
		promhttp.Handler().ServeHTTP(w, params.HTTPRequest)
	})
}

