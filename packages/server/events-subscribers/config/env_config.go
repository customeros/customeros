package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Logger              logger.Config
	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig
	Neo4j               config.Neo4jConfig
	Jaeger              tracing.JaegerConfig
	GrpcClientConfig    config.GrpcClientConfig

	CommonConfig config.CommonConfig
}
