package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/metrics"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	ApiPort     string `env:"PORT" envDefault:"10000" validate:"required"`
	MetricsPort string `env:"PORT_METRICS" envDefault:"10000" validate:"required"`
	Logger      logger.Config
	GraphQL     struct {
		PlaygroundEnabled    bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
		FixedComplexityLimit int  `env:"GRAPHQL_FIXED_COMPLEXITY_LIMIT" envDefault:"200"`
	}
	Admin struct {
		Key string `env:"ADMIN_KEY,required"`
	}
	Service struct {
		EventsProcessingPlatformUrl    string `env:"EVENTS_PROCESSING_PLATFORM_URL" validate:"required"`
		EventsProcessingPlatformApiKey string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" validate:"required"`
	}
	Postgres config.PostgresConfig
	Neo4j    config.Neo4jConfig
	Jaeger   tracing.Config
	Metrics  metrics.Config
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
