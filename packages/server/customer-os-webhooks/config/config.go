package config

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/joho/godotenv"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/metrics"
)

// TODO implement same approach as in packages/server/customer-os-api/config/config.go
type Config struct {
	App    App
	Common commonconf.CommonConfig
}

type CommonConfig struct {
	Logger         logger.Config
	Jaeger         tracing.JaegerConfig
	RabbitMQConfig commonconf.RabbitMQConfig
	Novu           commonconf.NovuConfig
	Temporal       commonconf.TemporalConfig
}

type App struct {
	ApiPort           string `env:"PORT" envDefault:"10004" validate:"required"`
	MetricsPort       string `env:"PORT_METRICS" envDefault:"10004" validate:"required"`
	ConcurrencyConfig ConcurrencyConfig
	Metrics           metrics.Config
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
		log.Fatalf("Error loading app configuration: %+v", err)
	}

	return &cfg, nil
}
