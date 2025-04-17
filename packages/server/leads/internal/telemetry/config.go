package telemetry

type OpenTelemetryConfig struct {
	Enabled     bool   `env:"OTEL_ENABLED" envDefault:"true"`
	Endpoint    string `env:"OTEL_ENDPOINT" envDefault:"otel-collector:4319"`
	ServiceName string `env:"OTEL_SERVICE_NAME" envDefault:"mailstack"`
	Timeout     int    `env:"OTEL_TIMEOUT" envDefault:"30"`
}
