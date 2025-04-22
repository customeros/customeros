package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	commonconf "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig           *AppConfig
	NATSConfig          *NATSConfig
	DataWarehouseConfig *DataWarehouseConfig
	CommonConfig        *commonconf.CommonConfig
}

type CommonConfig struct {
	Logger        logger.Config
	OpenTelemetry telemetry.OpenTelemetryConfig
	Postgres      commonconf.PostgresConfig
	Neo4j         commonconf.Neo4jConfig
}

func InitConfig() (*Config, error) {
	commonCfg := &CommonConfig{}

	config := &Config{
		AppConfig:           &AppConfig{},
		NATSConfig:          &NATSConfig{},
		DataWarehouseConfig: &DataWarehouseConfig{},
		CommonConfig: &commonconf.CommonConfig{
			Infrastructure: commonconf.InfrastructureConfig{
				LoggerConfig:        commonCfg.Logger,
				OpenTelemetryConfig: commonCfg.OpenTelemetry,
				PostgresConfig:      commonCfg.Postgres,
				Neo4jConfig:         commonCfg.Neo4j,
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
