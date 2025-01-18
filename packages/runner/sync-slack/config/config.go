package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronConfig "github.com/customeros/customeros/packages/runner/sync-slack/cron/config"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	CommonConfig commonConfig.CommonConfig

	Logger             logger.Config
	Jaeger             tracing.JaegerConfig
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
