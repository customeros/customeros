package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort string `env:"PORT"`

	ReacherConfig  ReacherConfig
	SmartyConfig   SmartyConfig
	IpDataConfig   IpDataConfig
	IpHunterConfig IpHunterConfig

	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
	Jaeger   tracing.JaegerConfig
	Logger   logger.Config
}

type ReacherConfig struct {
	ApiPath string `env:"REACHER_API_PATH" validate:"required"`
	Secret  string `env:"REACHER_SECRET" validate:"required"`
}

type SmartyConfig struct {
	AuthId    string `env:"SMARTY_AUTH_ID" validate:"required"`
	AuthToken string `env:"SMARTY_AUTH_TOKEN" validate:"required"`
}

type IpDataConfig struct {
	ApiUrl                   string `env:"IPDATA_API_URL" validate:"required"`
	ApiKey                   string `env:"IPDATA_API_KEY" validate:"required"`
	InvalidateCacheAfterDays int    `env:"IPDATA_INVALIDATE_CACHE_AFTER_IN_DAYS" default:"90"  validate:"required"`
}

type IpHunterConfig struct {
	ApiKey string `env:"IPHUNTER_API_KEY" validate:"required"`
}
