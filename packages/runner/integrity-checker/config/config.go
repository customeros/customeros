package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconfig "github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/cron/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/tracing"
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
	Logger logger.Config
	Jaeger tracing.Config
	Cron   cronconfig.Config

	AWS struct {
		Bucket                               string `env:"AWS_S3_BUCKET,required"`
		CloudWatchNamespace                  string `env:"AWS_CLOUDWATCH_METRICS_NAMESPACE,required" envDefault:"Openline"`
		MetricsDimensionEnvironment          string `env:"AWS_CLOUDWATCH_METRICS_DIMENSION_ENVIRONMENT,required" envDefault:"openline-dev"`
		MetricsDimensionNeo4jIntegrityChecks string `env:"AWS_CLOUDWATCH_METRICS_DIMENSION_NEO4J_INTEGRITY_CHECKS,required" envDefault:"Neo4jIntegrityChecks"`
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
