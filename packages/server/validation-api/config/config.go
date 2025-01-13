package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort string `env:"PORT"`

	SmartyConfig    SmartyConfig
	IpDataConfig    IpDataConfig
	IpHunterConfig  IpHunterConfig
	EmailConfig     EmailConfig
	ScrubbyIoConfig ScrubbyIoConfig
	TrueinboxConfig TrueinboxConfig
	EnrowConfig     EnrowConfig

	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig
	Neo4j               config.Neo4jConfig
	Jaeger              tracing.JaegerConfig
	Logger              logger.Config
}
