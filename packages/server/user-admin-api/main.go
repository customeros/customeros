package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/routes"
	"log"
)

func loadConfiguration() config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return cfg
}

func main() {
	config := loadConfiguration()

	graphqlClient := graphql.NewClient(config.CustomerOS.CustomerOsAPI)
	routes.Run(&config, graphqlClient)
}
