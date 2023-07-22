package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cron_config "github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/cron/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-dedup/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"log"
)

type Config struct {
	Neo4jDb struct {
		Target                string `env:"NEO4J_TARGET,required"`
		User                  string `env:"NEO4J_AUTH_USER,required,unset"`
		Pwd                   string `env:"NEO4J_AUTH_PWD,required,unset"`
		Realm                 string `env:"NEO4J_AUTH_REALM"`
		MaxConnectionPoolSize int    `env:"NEO4J_MAX_CONN_POOL_SIZE" envDefault:"100"`
		LogLevel              string `env:"NEO4J_LOG_LEVEL" envDefault:"WARNING"`
	}
	Logger        logger.Config
	Jaeger        tracing.Config
	Cron          cron_config.Config
	Organizations struct {
		AtLeastPerTenant                        int     `env:"ORGANIZATIONS_AT_LEAST_PER_TENANT" envDefault:"10"`
		ApiPageSize                             int     `env:"ORGANIZATIONS_API_PAGE_SIZE" envDefault:"100"`
		CompareWindowSize                       int     `env:"ORGANIZATIONS_COMPARE_WINDOW_SIZE" envDefault:"110"`
		MinConfidenceLvlDuplicateCheckByName    float64 `env:"ORGANIZATIONS_MIN_CONFIDENCE_LVL_DUPLICATE_CHECK_BY_NAME" envDefault:"0.7"`
		MinConfidenceLvlDuplicateCheckByDetails float64 `env:"ORGANIZATIONS_MIN_CONFIDENCE_LVL_DUPLICATE_CHECK_BY_DETAILS" envDefault:"0.8"`
	}
	Service struct {
		CustomerOsAdminAPI    string `env:"CUSTOMER_OS_ADMIN_API,required"`
		CustomerOsAdminAPIKey string `env:"CUSTOMER_OS_ADMIN_API_KEY,required"`
	}
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
