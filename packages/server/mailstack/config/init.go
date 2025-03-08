package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig               *AppConfig
	MailstackDatabaseConfig *MailstackDatabaseConfig
}

func InitConfig() (*Config, error) {
	config := &Config{
		AppConfig:               &AppConfig{},
		MailstackDatabaseConfig: &MailstackDatabaseConfig{},
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
