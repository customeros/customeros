package config

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"

type Config struct {
	PostgresConfig      config.PostgresConfig
	PostgresAsyncConfig config.PostgresAsyncConfig

	ApiPort  string `env:"PORT" envDefault:"10101" validate:"required"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`

	Anthropic struct {
		ApiPath string `env:"ANTHROPIC_API_PATH,required" envDefault:"https://api.anthropic.com/v1/messages"`
		ApiKey  string `env:"ANTHROPIC_API_KEY,required" envDefault:""`
	}
}
