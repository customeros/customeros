package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronConfig "github.com/openline-ai/openline-customer-os/packages/runner/sync-slack/cron/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Postgres           commonConfig.PostgresConfig
	Neo4j              commonConfig.Neo4jConfig
	Logger             logger.Config
	Jaeger             tracing.Config
	Cron               cronConfig.Config
	RawDataStoreDBName string `env:"RAW_DATA_STORE_DB_NAME,required" envDefault:"destination"`
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
