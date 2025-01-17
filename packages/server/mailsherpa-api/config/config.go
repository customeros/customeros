package config

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ApiPort string `env:"PORT"`

	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig
	Jaeger              tracing.JaegerConfig
	Logger              logger.Config
	EmailConfig         EmailConfig
}

type EmailConfig struct {
	EmailValidationFromDomain                  string `env:"EMAIL_VALIDATION_FROM_DOMAIN"`
	EmailDomainValidationCacheTtlDays          int    `env:"EMAIL_VALIDATION_DOMAIN_CACHE_TTL_DAYS" envDefault:"90" validate:"required"`
	EmailValidationCacheTtlDays                int    `env:"EMAIL_VALIDATION_CACHE_TTL_DAYS" envDefault:"14" validate:"required"`
	EmailValidationSkipProvidersCommaSeparated string `env:"EMAIL_VALIDATION_SKIP_PROVIDERS" envDefault:""`
}
