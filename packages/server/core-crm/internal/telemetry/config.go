package telemetry

type OpenTelemetryConfig struct {
	Enabled     bool   `env:"OTEL_ENABLED" envDefault:"true"`
	Endpoint    string `env:"OTEL_ENDPOINT" envDefault:"otel-collector:4317"`
	ServiceName string `env:"OTEL_SERVICE_NAME" envDefault:"core-crm"`
	Timeout     int    `env:"OTEL_TIMEOUT" envDefault:"30"`
}
