package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/joho/godotenv"

	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
)

type Config struct {
	AppConfig           *AppConfig
	Telemetry           *telemetry.OpenTelemetryConfig
	NATSConfig          *NATSConfig
	DataWarehouseConfig *DataWarehouseConfig
	CommonConfig        *commonconf.CommonConfig
}

type CommonConfig struct {
	Logger   logger.Config
	Postgres commonconf.PostgresConfig
	Neo4j    commonconf.Neo4jConfig
}

func InitConfig() (*Config, error) {
	commonCfg := &CommonConfig{}

	config := &Config{
		AppConfig:           &AppConfig{},
		Telemetry:           &telemetry.OpenTelemetryConfig{},
		NATSConfig:          &NATSConfig{},
		DataWarehouseConfig: &DataWarehouseConfig{},
		CommonConfig: &commonconf.CommonConfig{
			Infrastructure: commonconf.InfrastructureConfig{
				LoggerConfig:   commonCfg.Logger,
				PostgresConfig: commonCfg.Postgres,
				Neo4jConfig:    commonCfg.Neo4j,
			},
		},
	}

	err := godotenv.Load()
	if err != nil {
		log.Print("Unable to load .env file")
	}
	if err := env.Parse(config); err != nil {
		log.Fatalf("%+v", err)
	}
	if err := env.Parse(commonCfg); err != nil {
		log.Fatalf("%+v", err)
	}

	return config, nil
}
