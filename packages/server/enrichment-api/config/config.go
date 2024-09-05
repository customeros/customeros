package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/metrics"
	"log"
)

type Config struct {
	ApiPort             string `env:"PORT" envDefault:"10007" validate:"required"`
	MetricsPort         string `env:"PORT_METRICS" envDefault:"10007" validate:"required"`
	Logger              logger.Config
	Postgres            config.PostgresConfig
	Neo4j               config.Neo4jConfig
	Jaeger              tracing.JaegerConfig
	Metrics             metrics.Config
	ScrapinConfig       ScrapinConfig
	BetterContactConfig BetterContactConfig
	BrandfetchConfig    BrandfetchConfig
}

type ScrapinConfig struct {
	Url     string `env:"SCRAPIN_API_URL" envDefault:"https://api.scrapin.io" required:"true"`
	ApiKey  string `env:"SCRAPIN_API_KEY" required:"true"`
	TtlDays int    `env:"SCRAPIN_TTL_DAYS" envDefault:"90" required:"true"`
}

type BrandfetchConfig struct {
	Url     string `env:"BRANDFETCH_API_URL"`
	Limit   int    `env:"BRANDFETCH_LIMIT" envDefault:"250"`
	TtlDays int    `env:"BRANDFETCH_TTL_DAYS" envDefault:"180" required:"true"`
}

type BetterContactConfig struct {
	Url         string `env:"BETTER_CONTACT_API_URL" required:"true"`
	ApiKey      string `env:"BETTER_CONTACT_API_KEY" required:"true"`
	CallbackUrl string `env:"BETTER_CONTACT_CALLBACK_URL" required:"true"`
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
