package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/metrics"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	ApiPort     string `env:"PORT" envDefault:"10000" validate:"required"`
	MetricsPort string `env:"PORT_METRICS" envDefault:"10000" validate:"required"`
	Logger      logger.Config
	GraphQL     struct {
		PlaygroundEnabled    bool `env:"GRAPHQL_PLAYGROUND_ENABLED" envDefault:"false"`
		FixedComplexityLimit int  `env:"GRAPHQL_FIXED_COMPLEXITY_LIMIT" envDefault:"200"`
	}
	Admin struct {
		Key string `env:"ADMIN_KEY,required"`
	}
	Postgres struct {
		Host            string `env:"POSTGRES_HOST,required"`
		Port            string `env:"POSTGRES_PORT,required"`
		User            string `env:"POSTGRES_USER,required,unset"`
		Db              string `env:"POSTGRES_DB,required"`
		Password        string `env:"POSTGRES_PASSWORD,required,unset"`
		MaxConn         int    `env:"POSTGRES_DB_MAX_CONN"`
		MaxIdleConn     int    `env:"POSTGRES_DB_MAX_IDLE_CONN"`
		ConnMaxLifetime int    `env:"POSTGRES_DB_CONN_MAX_LIFETIME"`
		LogLevel        string `env:"POSTGRES_LOG_LEVEL" envDefault:"WARN"`
	}
	Neo4j struct {
		Target                          string `env:"NEO4J_TARGET,required"`
		User                            string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                             string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                           string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize           int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		ConnectionAcquisitionTimeoutSec int    `env:"NEO4J_CONN_ACQUISITION_TIMEOUT_SEC" envDefault:"60"`
		LogLevel                        string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
	Service struct {
		EventsProcessingPlatformUrl    string `env:"EVENTS_PROCESSING_PLATFORM_URL" validate:"required"`
		EventsProcessingPlatformApiKey string `env:"EVENTS_PROCESSING_PLATFORM_API_KEY" validate:"required"`
	}
	Jaeger  tracing.Config
	Metrics metrics.Config
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
