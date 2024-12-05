package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to capture the status code
		rw := &statusRecordingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Proceed with the next handler
		next.ServeHTTP(rw, r)

		// Record metrics
		requestCount.With(prometheus.Labels{
			"path":   r.URL.Path,
			"method": r.Method,
			"status": strconv.Itoa(rw.statusCode),
		}).Inc()
	})
}

type statusRecordingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusRecordingResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
