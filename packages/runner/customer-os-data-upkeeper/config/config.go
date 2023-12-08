package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconf "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/cron/config"
	commconf "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Postgres         commconf.PostgresConfig
	Neo4j            commconf.Neo4jConfig
	Logger           logger.Config
	Jaeger           tracing.JaegerConfig
	Cron             cronconf.Config
	EventsProcessing EventsProcessingConfig
	ProcessConfig    ProcessConfig
}

type EventsProcessingConfig struct {
	EventsProcessingPlatformEnabled bool   `env:"EVENTS_PROCESSING_PLATFORM_ENABLED" envDefault:"true"`
	EventsProcessingPlatformUrl     string `env:"EVENTS_PROCESSING_PLATFORM_URL"`
	EventsProcessingPlatformApiKey  string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY"`
}

type ProcessConfig struct {
	WebScrapedOrganizationsPerCycle int `env:"WEB_SCRAPED_ORGANIZATIONS_PER_CYCLE" envDefault:"200"`
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
