package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	common_config "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/joho/godotenv"

	"github.com/customeros/customeros/packages/server/ai/internal/logger"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
)

type Config struct {
	AppConfig     *AppConfig
	Telemetry     *telemetry.OpenTelemetryConfig
	Logger        *logger.Config
	Postgres      common_config.PostgresConfig
	Neo4j         common_config.Neo4jConfig
	DataWarehouse common_config.DataWarehouseConfig
	NATS          *NATSConfig
	Anthropic     *AnthropicConfig
	Deepseek      *DeepseekConfig
	Groq          *GroqConfig
	Gemini        *GeminiConfig
}

func InitConfig() (*Config, error) {
	config := &Config{
		AppConfig:     &AppConfig{},
		Telemetry:     &telemetry.OpenTelemetryConfig{},
		Logger:        &logger.Config{},
		NATS:          &NATSConfig{},
		Anthropic:     &AnthropicConfig{},
		Deepseek:      &DeepseekConfig{},
		Groq:          &GroqConfig{},
		Gemini:        &GeminiConfig{},
		Postgres:      common_config.PostgresConfig{},
		Neo4j:         common_config.Neo4jConfig{},
		DataWarehouse: common_config.DataWarehouseConfig{},
	}

	err := godotenv.Load()
	if err != nil {
		log.Print("Unable to load .env file")
	}
	if err := env.Parse(config); err != nil {
		log.Fatalf("%+v", err)
	}

	return config, nil
}
