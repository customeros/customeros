package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	Logger              logger.Config
	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig
	Neo4j               config.Neo4jConfig
	Jaeger              tracing.JaegerConfig
	RabbitMQ            config.RabbitMQConfig
	GrpcClientConfig    config.GrpcClientConfig
	InternalServices    InternalServices
	ExternalServices    ExternalServices
}

type InternalServices struct {
	EnrichmentApi config.EnrichmentAPIConfig
	AiApi         config.AiAPIConfig
	ValidationApi config.ValidationAPIConfig
}

type ExternalServices struct {
	CloudflareConfig config.CloudflareConfig
	NamecheapConfig  config.NamecheapConfig
	OpenSRSConfig    config.OpenSRSConfig
	AnthropicConfig  config.AnthropicConfig
	NovuConfig       config.NovuConfig
}
