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

	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
	Jaeger   tracing.JaegerConfig
	Logger   logger.Config
}

type SmartyConfig struct {
	AuthId    string `env:"SMARTY_AUTH_ID" validate:"required"`
	AuthToken string `env:"SMARTY_AUTH_TOKEN" validate:"required"`
}

type IpDataConfig struct {
	ApiUrl             string `env:"IPDATA_API_URL" validate:"required"`
	ApiKey             string `env:"IPDATA_API_KEY" validate:"required"`
	IpDataCacheTtlDays int    `env:"IPDATA_CACHE_TTL_DAYS" envDefault:"90" validate:"required"`
}

type IpHunterConfig struct {
}

type EmailConfig struct {
	EmailValidationFromDomain                  string `env:"EMAIL_VALIDATION_FROM_DOMAIN"`
	EmailDomainValidationCacheTtlDays          int    `env:"EMAIL_VALIDATION_DOMAIN_CACHE_TTL_DAYS" envDefault:"90" validate:"required"`
	EmailValidationCacheTtlDays                int    `env:"EMAIL_VALIDATION_CACHE_TTL_DAYS" envDefault:"14" validate:"required"`
	EmailValidationSkipProvidersCommaSeparated string `env:"EMAIL_VALIDATION_SKIP_PROVIDERS" envDefault:""`
}

type ScrubbyIoConfig struct {
	ApiUrl               string `env:"SCRUBBY_IO_API_URL" envDefault:"https://api.scrubby.io" validate:"required"`
	ApiKey               string `env:"SCRUBBY_IO_API_KEY" validate:"required"`
	CacheTtlDays         int    `env:"SCRUBBY_IO_CACHE_TTL_DAYS" envDefault:"90" validate:"required"`
	ScrubbyIoCallbackUrl string `env:"SCRUBBY_IO_CALLBACK_URL"`
}
