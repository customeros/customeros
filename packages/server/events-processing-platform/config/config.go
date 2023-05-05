package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstroredb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
)

var configPath string

func init() {
}

type GRPC struct {
	Port        string `env:"GRPC_PORT" envDefault:":5001" validate:"required"`
	Development bool   `env:"GRPC_DEVELOPMENT" envDefault:"false"`
}

type Subscriptions struct {
	PoolSize                    int    `env:"EVENT_STORE_SUBSCRIPTIONS_POOL_SIZE" envDefault:"60" validate:"required,gte=0"`
	PhoneNumberPrefix           string `env:"EVENT_STORE_SUBSCRIPTIONS_PHONE_PREFIX" envDefault:"phone_number-" validate:"required"`
	EmailPrefix                 string `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_PREFIX" envDefault:"email-" validate:"required"`
	UserPrefix                  string `env:"EVENT_STORE_SUBSCRIPTIONS_USER_PREFIX" envDefault:"user-" validate:"required"`
	ContactPrefix               string `env:"EVENT_STORE_SUBSCRIPTIONS_CONTACT_PREFIX" envDefault:"contact-" validate:"required"`
	OrganizationPrefix          string `env:"EVENT_STORE_SUBSCRIPTIONS_ORGANIZATION_PREFIX" envDefault:"organization-" validate:"required"`
	GraphSubscription           GraphSubscription
	EmailValidationSubscription EmailValidationSubscription
}

type GraphSubscription struct {
	Enabled   bool   `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_ENABLED" envDefault:"true"`
	GroupName string `env:"EVENT_STORE_SUBSCRIPTIONS_GRAPH_GROUP_NAME" envDefault:"graph-v1" validate:"required"`
}

type EmailValidationSubscription struct {
	Enabled   bool   `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_ENABLED" envDefault:"false"`
	GroupName string `env:"EVENT_STORE_SUBSCRIPTIONS_EMAIL_VALIDATION_GROUP_NAME" envDefault:"emailValidation-v1" validate:"required"`
}

type Neo4j struct {
	Target                          string `env:"NEO4J_TARGET" validate:"required"`
	User                            string `env:"NEO4J_AUTH_USER,unset" validate:"required"`
	Pwd                             string `env:"NEO4J_AUTH_PWD,unset" validate:"required"`
	Realm                           string `env:"NEO4J_AUTH_REALM"`
	MaxConnectionPoolSize           int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
	ConnectionAcquisitionTimeoutSec int    `env:"NEO4J_CONN_ACQUISITION_TIMEOUT_SEC" envDefault:"60"`
	LogLevel                        string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
}

type Services struct {
	ValidationApi    string `env:"VALIDATION_API" validate:"required"`
	ValidationApiKey string `env:"VALIDATION_API_KEY" validate:"required"`
}

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"events-processing-platform"`
	Logger           logger.Config
	GRPC             GRPC
	EventStoreConfig eventstroredb.EventStoreConfig
	Subscriptions    Subscriptions
	Neo4j            Neo4j
	Jaeger           tracing.Config
	Services         Services
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
