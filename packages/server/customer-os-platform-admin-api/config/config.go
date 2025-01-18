package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/validator"
	"github.com/customeros/customeros/packages/server/customer-os-platform-admin-api/metrics"
	"log"
)

type Config struct {
	App    App
	Common config.CommonConfig
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

	err := validator.GetValidator().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
