// Package metrics provides Prometheus metrics collection for the API.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal    *prometheus.CounterVec
	HTTPRequestDuration  *prometheus.HistogramVec
	HTTPRequestsInFlight prometheus.Gauge
	HTTPResponseSize     *prometheus.HistogramVec

	// Database metrics
	DBQueriesTotal   *prometheus.CounterVec
	DBQueryDuration  *prometheus.HistogramVec
	DBConnectionPool *prometheus.GaugeVec

	// Auth metrics
	AuthAttemptsTotal *prometheus.CounterVec
	TokensIssued      *prometheus.CounterVec

	// Business metrics
	UsersTotal    prometheus.Gauge
	MessagesTotal prometheus.Gauge
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		HTTPRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "http_requests_in_flight",
				Help:      "Number of HTTP requests currently being processed",
			},
		),
		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		DBQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "db_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "status"},
		),
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "db_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		DBConnectionPool: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_connection_pool",
				Help:      "Database connection pool statistics",
			},
			[]string{"state"},
		),
		AuthAttemptsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "auth_attempts_total",
				Help:      "Total number of authentication attempts",
			},
			[]string{"method", "status"},
		),
		TokensIssued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tokens_issued_total",
				Help:      "Total number of tokens issued",
			},
			[]string{"type"},
		),
		UsersTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "users_total",
				Help:      "Total number of users in the system",
			},
		),
		MessagesTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "messages_total",
				Help:      "Total number of messages in the system",
			},
		),
	}
}

// DefaultMetrics is the default metrics instance
var DefaultMetrics = NewMetrics("subspace")

// RecordHTTPRequest records metrics for an HTTP request
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration float64, size int) {
	m.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	m.HTTPResponseSize.WithLabelValues(method, path).Observe(float64(size))
}

// RecordDBQuery records metrics for a database query
func (m *Metrics) RecordDBQuery(operation, table, status string, duration float64) {
	m.DBQueriesTotal.WithLabelValues(operation, table, status).Inc()
	m.DBQueryDuration.WithLabelValues(operation, table).Observe(duration)
}

// RecordAuthAttempt records an authentication attempt
func (m *Metrics) RecordAuthAttempt(method, status string) {
	m.AuthAttemptsTotal.WithLabelValues(method, status).Inc()
}

// RecordTokenIssued records a token issuance
func (m *Metrics) RecordTokenIssued(tokenType string) {
	m.TokensIssued.WithLabelValues(tokenType).Inc()
}

// UpdateDBPoolStats updates database connection pool statistics
func (m *Metrics) UpdateDBPoolStats(open, inUse, idle int) {
	m.DBConnectionPool.WithLabelValues("open").Set(float64(open))
	m.DBConnectionPool.WithLabelValues("in_use").Set(float64(inUse))
	m.DBConnectionPool.WithLabelValues("idle").Set(float64(idle))
}

// UpdateBusinessMetrics updates business-related metrics
func (m *Metrics) UpdateBusinessMetrics(users, messages int) {
	m.UsersTotal.Set(float64(users))
	m.MessagesTotal.Set(float64(messages))
}
