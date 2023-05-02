package main

import (
	"encoding/json"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"log"
	"os"
)

func main() {
	config := loadConfiguration()

	graphqlClient := graphql.NewClient(config.Service.CustomerOsAPI)
	services := service.InitServices(graphqlClient, &config)
	hub := ContactHub.NewContactHub()
	go hub.Run()
	routes.Run(&config, hub, services) // run this as a background goroutine

}

type ServiceConfig struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func loadConfiguration() commsApiConfig.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := commsApiConfig.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	file, err := os.Open("torrey-test-email-service-ed4fb89333e7.json")
	if err != nil {
		log.Println("[WARNING] Error loading .json file")
		return cfg
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var serviceConfig ServiceConfig
	err = decoder.Decode(&serviceConfig)
	if err != nil {
		log.Println("[WARNING] Error decoding .json file")
		return cfg
	}

	cfg.GMail.ServiceEmailAddress = serviceConfig.ClientEmail
	cfg.GMail.ServicePrivateKey = serviceConfig.PrivateKey
	return cfg
}
