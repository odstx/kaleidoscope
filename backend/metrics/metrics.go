package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed",
		},
		[]string{"method", "path"},
	)

	UserOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_operations_total",
			Help: "Total number of user operations",
		},
		[]string{"operation", "status"},
	)

	UserOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_operation_duration_seconds",
			Help:    "Duration of user operations in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)

	ActiveUsersGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users_total",
			Help: "Total number of registered users",
		},
	)

	TOTPOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "totp_operations_total",
			Help: "Total number of TOTP operations",
		},
		[]string{"operation", "status"},
	)

	AuthOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_operations_total",
			Help: "Total number of authentication operations",
		},
		[]string{"type", "status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"operation", "table"},
	)

	EtcdOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "etcd_operation_duration_seconds",
			Help:    "Duration of etcd operations in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"operation"},
	)

	EtcdOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "etcd_operations_total",
			Help: "Total number of etcd operations",
		},
		[]string{"operation", "status"},
	)

	MicroserviceRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "microservice_requests_total",
			Help: "Total number of microservice proxy requests",
		},
		[]string{"app", "status"},
	)

	MicroserviceRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "microservice_request_duration_seconds",
			Help:    "Duration of microservice proxy requests in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, 30},
		},
		[]string{"app"},
	)

	MicroserviceCircuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "microservice_circuit_breaker_state",
			Help: "Circuit breaker state for microservice (0=closed, 1=open, 2=half-open)",
		},
		[]string{"app"},
	)

	MicroserviceInstanceHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "microservice_instance_health",
			Help: "Health status of microservice instances (1=healthy, 0=unhealthy)",
		},
		[]string{"app", "instance_id"},
	)

	AgentRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "agent_requests_total",
			Help: "Total number of agent chat requests",
		},
		[]string{"status"},
	)

	AgentRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "agent_request_duration_seconds",
			Help:    "Duration of agent chat requests in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60},
		},
		[]string{},
	)

	EmailTasksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "email_tasks_total",
			Help: "Total number of email tasks enqueued",
		},
		[]string{"type", "status"},
	)

	SecurityEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "security_events_total",
			Help: "Total number of security events",
		},
		[]string{"event_type", "severity"},
	)

	SecurityEventDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "security_event_duration_seconds",
			Help:    "Duration of security event processing",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"event_type"},
	)
)

func RecordUserOperation(operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	UserOperationsTotal.WithLabelValues(operation, status).Inc()
	UserOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

func RecordTOTPOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	TOTPOperationsTotal.WithLabelValues(operation, status).Inc()
}

func RecordAuthOperation(authType string, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	AuthOperationsTotal.WithLabelValues(authType, status).Inc()
}

func RecordDatabaseQuery(operation, table string, duration time.Duration) {
	DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func RecordEtcdOperation(operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	EtcdOperationsTotal.WithLabelValues(operation, status).Inc()
	EtcdOperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

func RecordMicroserviceRequest(app string, duration time.Duration, statusCode int) {
	status := strconv.Itoa(statusCode)
	MicroserviceRequestsTotal.WithLabelValues(app, status).Inc()
	MicroserviceRequestDuration.WithLabelValues(app).Observe(duration.Seconds())
}

func RecordCircuitBreakerState(app string, state float64) {
	MicroserviceCircuitBreakerState.WithLabelValues(app).Set(state)
}

func RecordInstanceHealth(app, instanceID string, healthy bool) {
	value := 0.0
	if healthy {
		value = 1.0
	}
	MicroserviceInstanceHealth.WithLabelValues(app, instanceID).Set(value)
}

func RecordAgentRequest(duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "failure"
	}
	AgentRequestsTotal.WithLabelValues(status).Inc()
	AgentRequestDuration.WithLabelValues().Observe(duration.Seconds())
}

func RecordEmailTask(taskType string, success bool) {
	status := "enqueued"
	if !success {
		status = "failed"
	}
	EmailTasksTotal.WithLabelValues(taskType, status).Inc()
}

func SetActiveUsers(count float64) {
	ActiveUsersGauge.Set(count)
}

func RecordSecurityEvent(eventType, severity string, duration time.Duration) {
	SecurityEventsTotal.WithLabelValues(eventType, severity).Inc()
	SecurityEventDuration.WithLabelValues(eventType).Observe(duration.Seconds())
}