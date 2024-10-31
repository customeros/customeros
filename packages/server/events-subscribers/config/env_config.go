package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Logger           logger.Config
	Postgres         config.PostgresConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
	RabbitMQ         config.RabbitMQConfig
	NovuConfig       config.NovuConfig
	GrpcClientConfig config.GrpcClientConfig
	InternalServices InternalServices
}

type InternalServices struct {
	EnrichmentApi config.EnrichmentAPIConfig
	AiApi         config.AiAPIConfig
}
