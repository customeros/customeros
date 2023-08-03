package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	cronconfig "github.com/openline-ai/openline-customer-os/packages/runner/integrity-checker/cron/config"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"log"
)

type Config struct {
	Neo4j  commonConfig.Neo4jConfig
	Logger logger.Config
	Jaeger tracing.Config
	Cron   cronconfig.Config

	AWS struct {
		Bucket                               string `env:"AWS_S3_BUCKET,required"`
		Region                               string `env:"AWS_REGION,required"`
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
