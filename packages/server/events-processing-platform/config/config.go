package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstoredb"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger           logger.Config
	EventStoreConfig eventstoredb.EventStoreConfig
	Neo4j            config.Neo4jConfig
	Postgres         config.PostgresConfig
	Jaeger           tracing.JaegerConfig
	GRPC             GRPC
	Services         Services
	Utils            Utils
}

type GRPC struct {
	Port        string `env:"GRPC_PORT" envDefault:":5001" validate:"required"`
	Development bool   `env:"GRPC_DEVELOPMENT" envDefault:"false"`
	ApiKey      string `env:"GRPC_API_KEY" validate:"required"`
}

type Services struct {
	FileStoreApiConfig fsc.FileStoreApiConfig
}

type Utils struct {
	RetriesOnOptimisticLockException int `env:"UTILS_RETRIES_ON_OPTIMISTIC_LOCK" envDefault:"5"`
}

func InitConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	if err = env.Parse(&cfg); err != nil {
		return nil, err
	}

	err = validator.GetValidator().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
