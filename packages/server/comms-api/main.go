package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/chatHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	config := loadConfiguration()

	mh := chatHub.NewHub()
	go mh.Run()

	graphqlClient := graphql.NewClient(config.Service.CustomerOsAPI)
	customerOSService := service.NewCustomerOSService(graphqlClient, &config)
	// Our server will live in the routes package
	go routes.Run(&config, mh, customerOSService) // run this as a background goroutine

	// Initialize the generated User service.
	//df := util.MakeDialFactory(&config)
	//svc := service.NewSendMessageService(&config, df, mh)

	log.Printf("Attempting to start GRPC server")
	// Create a new gRPC server (you can wire multiple services to a single server).
	server := grpc.NewServer()

	// Register the MessageDeprecate Item service with the server.
	//proto.RegisterMessageEventServiceServer(server, svc)

	// Open port for listening to traffic.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Service.GRPCPort))
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	} else {
		log.Printf("server started on: %s", fmt.Sprintf(":%d", config.Service.GRPCPort))
	}

	// Listen for traffic indefinitely.
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}

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
