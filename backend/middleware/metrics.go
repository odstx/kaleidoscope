package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"kaleidoscope/metrics"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
		[]string{"method", "path"},
	)
)

func RecordUserOperation(operation string, duration time.Duration, success bool) {
	metrics.RecordUserOperation(operation, duration, success)
}

func RecordTOTPOperation(operation string, success bool) {
	metrics.RecordTOTPOperation(operation, success)
}

func RecordAuthOperation(authType string, success bool) {
	metrics.RecordAuthOperation(authType, success)
}

func RecordDatabaseQuery(operation, table string, duration time.Duration) {
	metrics.RecordDatabaseQuery(operation, table, duration)
}

func RecordEtcdOperation(operation string, duration time.Duration, success bool) {
	metrics.RecordEtcdOperation(operation, duration, success)
}

func RecordMicroserviceRequest(app string, duration time.Duration, statusCode int) {
	metrics.RecordMicroserviceRequest(app, duration, statusCode)
}

func RecordCircuitBreakerState(app string, state float64) {
	metrics.RecordCircuitBreakerState(app, state)
}

func RecordInstanceHealth(app, instanceID string, healthy bool) {
	metrics.RecordInstanceHealth(app, instanceID, healthy)
}

func RecordAgentRequest(duration time.Duration, success bool) {
	metrics.RecordAgentRequest(duration, success)
}

func RecordEmailTask(taskType string, success bool) {
	metrics.RecordEmailTask(taskType, success)
}

func SetActiveUsers(count float64) {
	metrics.SetActiveUsers(count)
}

func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		httpRequestsInFlight.WithLabelValues(method, path).Inc()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		httpRequestsInFlight.WithLabelValues(method, path).Dec()
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}