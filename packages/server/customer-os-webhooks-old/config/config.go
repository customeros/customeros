package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/metrics"
	"log"
)

type Config struct {
	ApiPort           string `env:"PORT" envDefault:"10004" validate:"required"`
	MetricsPort       string `env:"PORT_METRICS" envDefault:"10004" validate:"required"`
	GrpcClientConfig  config.GrpcClientConfig
	ConcurrencyConfig ConcurrencyConfig
	Logger            logger.Config
	Postgres          config.PostgresConfig
	Neo4j             config.Neo4jConfig
	RabbitMQConfig    config.RabbitMQConfig
	Jaeger            tracing.JaegerConfig
	Metrics           metrics.Config

	BetterContactCallbackApiKey string `env:"BETTER_CONTACT_CALLBACK_API_KEY" validate:"required"`
	EnrowCallbackApiKey         string `env:"ENROW_CALLBACK_API_KEY" validate:"required"`

	Slack struct {
		NotifyPostmarkEmail    string `env:"SLACK_NOTIFY_POSTMARK_EMAIL"`
		NotifyFlowGoalAchieved string `env:"SLACK_NOTIFY_FLOW_GOAL_ACHIEVED"`
	}
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
