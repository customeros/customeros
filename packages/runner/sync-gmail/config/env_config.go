package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Neo4jDb config.Neo4jConfig

	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig

	RabbitMQConfig   config.RabbitMQConfig
	GrpcClientConfig config.GrpcClientConfig

	CommonConfig config.CommonConfig

	SyncData struct {
		CronSync string `env:"CRON_SYNC" envDefault:"0 */1 * * * *"`
	}

	Jaeger tracing.JaegerConfig
	Logger logger.Config
}

const MAX_EMAILS_PER_RUN = 50
