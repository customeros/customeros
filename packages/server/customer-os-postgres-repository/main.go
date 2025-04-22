package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	common_config "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/joho/godotenv"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
)

type Config struct {
	PostgresConfig      common_config.PostgresConfig
	PostgresAsyncConfig common_config.PostgresAsyncConfig
}

func main() {
	cmnConfig := loadConfig()

	postgresDb, err := common_config.InitPostgres(&cmnConfig)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	repo := postgres_repository.InitRepositories(postgresDb)
	if err = repo.AutoMigrate(postgresDb); err != nil {
		log.Fatalf("failed to run auto-migration: %v", err)
	}

	log.Print("Auto-migration completed successfully")
}

func loadConfig() common_config.CommonConfig {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found. Proceeding with system environment variables. Error: %v", err)
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error loading app configuration: %+v", err)
	}

	return common_config.CommonConfig{
		Infrastructure: common_config.InfrastructureConfig{
			PostgresConfig: cfg.PostgresConfig,
		},
	}
}
