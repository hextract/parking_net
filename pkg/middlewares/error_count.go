package middlewares

import (
	"fmt"
	"net/http"
)

// responseWriterWrapper - для захвата статуса ответа
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// PrometheusErrorMiddleware - отслеживает ошибки
func (metrics *PrometheusMetrics) PrometheusErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.URL.Path

		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)

		if wrapper.statusCode >= 400 {
			statusText := http.StatusText(wrapper.statusCode)
			if statusText == "" {
				statusText = fmt.Sprintf("%d", wrapper.statusCode)
			}
			metrics.ErrorResponses.WithLabelValues(endpoint, fmt.Sprintf("%d", wrapper.statusCode), statusText).Inc()
		}
	})
}
