package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
)

type Config struct {
	ApiPort string `env:"PORT"`
	AppKey  string `env:"APP_KEY"`

	PostgresConfig config.PostgresConfig
	Neo4jConfig    config.Neo4jConfig
	Jaeger         telemetry.JaegerConfig
	OpenTelemetry  telemetry.OpenTelemetryConfig
	Logger         logger.Config
}
