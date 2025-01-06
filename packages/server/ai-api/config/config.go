package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/metrics"
)

type Config struct {
	Logger           logger.Config
	Postgres         config.PostgresConfig
	Neo4j            config.Neo4jConfig
	Jaeger           tracing.JaegerConfig
	Metrics          metrics.Config
	ApiPort          string `env:"PORT" envDefault:"10016" validate:"required"`
	MetricsPort      string `env:"PORT_METRICS" envDefault:"10016" validate:"required"`
	AnthropicApiKey  string `env:"ANTHROPIC_API_KEY" required:"true"`
	AnthropicApiPath string `env:"ANTHROPIC_API_URL" envDefault:"https://api.anthropic.com/v1/messages" required:"true"`
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}

	err := validator.GetValidator().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
