package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconf "github.com/openline-ai/openline-customer-os/packages/runner/sync-tracking/cron/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Neo4j            config.Neo4jConfig
	Postgres         config.PostgresConfig
	GrpcClientConfig config.GrpcClientConfig
	Logger           logger.Config
	Jaeger           tracing.JaegerConfig
	RabbitMQConfig   config.RabbitMQConfig

	Cron cronconf.Config

	SnitcherApi struct {
		Url    string `env:"SNITCHER_API_URL"`
		ApiKey string `env:"SNITCHER_API_KEY"`
	}
	ValidationApi struct {
		Url    string `env:"VALIDATION_API" validate:"required"`
		ApiKey string `env:"VALIDATION_API_KEY" validate:"required"`
	}

	SlackBotApiKey string `env:"SLACK_BOT_API_KEY"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("Failed loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	return &cfg
}
