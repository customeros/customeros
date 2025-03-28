package telemetry

type OpenTelemetryConfig struct {
	Enabled     bool   `env:"OTEL_ENABLED" envDefault:"true"`
	Endpoint    string `env:"OTEL_ENDPOINT"`
	ServiceName string `env:"OTEL_SERVICE_NAME"`
	Timeout     int    `env:"OTEL_TIMEOUT"`
}
