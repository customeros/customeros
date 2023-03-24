package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/chatHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"log"
)

func main() {
	config := loadConfiguration()

	mh := chatHub.NewHub()
	go mh.Run()

	graphqlClient := graphql.NewClient(config.Service.CustomerOsAPI)
	customerOSService := service.NewCustomerOSService(graphqlClient, &config)
	mailService := service.NewMailService(&config, customerOSService)
	// Our server will live in the routes package
	routes.Run(&config, mh, customerOSService, mailService) // run this as a background goroutine

}

func loadConfiguration() commsApiConfig.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := commsApiConfig.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return cfg
}
