package metrics

type Config struct {
	PrometheusPath string `env:"PROMETHEUS_PATH" envDefault:"/metrics" validate:"required"`
}

// Prometheus metrics
var ()
