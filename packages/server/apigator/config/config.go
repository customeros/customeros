package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort   string `env:"PORT"`
	JwtSecret string `env:"JWT_SECRET"`

	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig
	Neo4jConfig         config.Neo4jConfig
	Jaeger              tracing.JaegerConfig
	Logger              logger.Config
}
