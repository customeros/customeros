package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/customeros/customeros/packages/server/scraper/internal/logger"
	"github.com/customeros/customeros/packages/server/scraper/internal/telemetry"
)

type Config struct {
	AppConfig             *AppConfig
	Logger                *logger.Config
	NATSConfig            *NATSConfig
	OpenTelemetry         *telemetry.OpenTelemetryConfig
	ScraperDatabeseConfig *ScraperDatabaseConfig
	DataWarehouseConfig   *DataWarehouseConfig
}

func InitConfig() (*Config, error) {
	config := &Config{
		AppConfig:             &AppConfig{},
		Logger:                &logger.Config{},
		NATSConfig:            &NATSConfig{},
		OpenTelemetry:         &telemetry.OpenTelemetryConfig{},
		ScraperDatabeseConfig: &ScraperDatabaseConfig{},
		DataWarehouseConfig:   &DataWarehouseConfig{},
	}

	err := godotenv.Load()
	if err != nil {
		log.Print("Unable to load .env file")
	}

	err = env.Parse(config)
	if err != nil {
		log.Fatalf("Error loading mailstack config: %v", err)
	}

	return config, nil
}
