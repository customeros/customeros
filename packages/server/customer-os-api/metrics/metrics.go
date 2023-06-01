package metrics

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	PrometheusPath string `env:"PROMETHEUS_PATH" envDefault:"/metrics" validate:"required"`
	PrometheusPort string `env:"PROMETHEUS_PORT" envDefault:"19000" validate:"required"`
}

// Prometheus metrics
var (
	MetricsGraphqlRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customer_os_api_graphql_requests_total",
			Help: "Total number of GraphQL requests",
		},
		[]string{"name", "status"},
	)

	MetricsGraphqlRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "customer_os_api_graphql_request_duration_seconds",
			Help:    "GraphQL request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"name"},
	)

	MetricsGraphqlRequestErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "customer_os_api_graphql_request_errors_total",
			Help: "Total number of GraphQL request errors",
		},
		[]string{"name"},
	)
)
