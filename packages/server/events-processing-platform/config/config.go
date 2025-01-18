package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/events/eventstoredb"
)

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger           logger.Config
	EventStoreConfig eventstoredb.EventStoreConfig
	Jaeger           tracing.JaegerConfig
	GRPC             GRPC
	Utils            Utils

	CommonServices config.CommonConfig
}

type GRPC struct {
	Port        string `env:"GRPC_PORT" envDefault:":5001" validate:"required"`
	Development bool   `env:"GRPC_DEVELOPMENT" envDefault:"false"`
	ApiKey      string `env:"GRPC_API_KEY" validate:"required"`
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
