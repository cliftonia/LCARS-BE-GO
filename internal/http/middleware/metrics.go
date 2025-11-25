// Package middleware provides HTTP middleware functions for the API.
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/metrics"
	"github.com/gorilla/mux"
)

// Metrics returns middleware that records Prometheus metrics for HTTP requests
func Metrics(m *metrics.Metrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Track in-flight requests
			m.HTTPRequestsInFlight.Inc()
			defer m.HTTPRequestsInFlight.Dec()

			// Wrap response writer to capture status and size
			// Reuse the existing responseWriter from logger.go
			rw := newResponseWriter(w)

			// Process request
			next.ServeHTTP(rw, r)

			// Get route pattern for consistent labeling
			path := getRoutePath(r)
			duration := time.Since(start).Seconds()
			status := strconv.Itoa(rw.statusCode)

			// Record metrics
			m.RecordHTTPRequest(r.Method, path, status, duration, int(rw.written))
		})
	}
}

// getRoutePath extracts the route pattern from the request
// Falls back to the actual path if no route is matched
func getRoutePath(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route != nil {
		if path, err := route.GetPathTemplate(); err == nil {
			return path
		}
	}
	return r.URL.Path
}
