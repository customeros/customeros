package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Neo4jDb    config.Neo4jConfig
	PostgresDb config.PostgresConfig

	GrpcClientConfig config.GrpcClientConfig

	SyncData struct {
		CronSync string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
	}

	Jaeger tracing.JaegerConfig
	Logger logger.Config
}
