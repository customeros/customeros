package config

import (
	"github.com/caarlos0/env/v6"
	commconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/customer-os-platform-admin-api/metrics"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	App    App
	Common *commconf.CommonConfig
}

type CommonConfig struct {
	Logger        logger.Config
	Jaeger        tracing.JaegerConfig
	Postgres      commconf.PostgresConfig
	PostgresAsync commconf.PostgresAsyncConfig
	Neo4j         commconf.Neo4jConfig
	GrpcClient    commconf.GrpcClientConfig
}

type App struct {
	ApiPort     string `env:"PORT" envDefault:"10005" validate:"required"`
	MetricsPort string `env:"PORT_METRICS" envDefault:"10005" validate:"required"`
	Logger      logger.Config
	Metrics     metrics.Config
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Print("Error loading .env file")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v", err)
	}
	err := validator.GetValidator().Struct(cfg.App)
	if err != nil {
		return nil, err
	}

	cmnCfg := CommonConfig{}
	if err := env.Parse(&cmnCfg); err != nil {
		log.Fatalf("%+v", err)
	}
	err = validator.GetValidator().Struct(cmnCfg)
	if err != nil {
		return nil, err
	}

	cfg.Common = &commconf.CommonConfig{
		Infrastructure: commconf.InfrastructureConfig{
			LoggerConfig:        cmnCfg.Logger,
			JaegerConfig:        cmnCfg.Jaeger,
			PostgresConfig:      cmnCfg.Postgres,
			PostgresAsyncConfig: cmnCfg.PostgresAsync,
			Neo4jConfig:         cmnCfg.Neo4j,
			GrpcClientConfig:    cmnCfg.GrpcClient,
		},
	}

	return &cfg, nil
}
