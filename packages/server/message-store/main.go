package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/lib/pq"
	"github.com/machinebox/graphql"
	c "github.com/openline-ai/openline-customer-os/packages/server/message-store/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	conf := c.Config{}
	if err := env.Parse(&conf); err != nil {
		fmt.Printf("missing required config")
		return
	}
	var connUrl = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Name, conf.DB.Password)
	log.Printf("Connecting to database %s", connUrl)
	client, err := gen.Open("postgres", connUrl)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	graphqlClient := graphql.NewClient(conf.Service.CustomerOsAPI)

	// Initialize the generated User service.
	svc := service.NewMessageService(client, graphqlClient)

	// Create a new gRPC server (you can wire multiple services to a single server).
	server := grpc.NewServer()

	// Register the Message Item service with the server.
	pb.RegisterMessageStoreServiceServer(server, svc)

	// Open port for listening to traffic.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Service.ServerPort))
	if err != nil {
		log.Fatalf("failed listening: %s", err)
	} else {
		log.Printf("server started on: %s", fmt.Sprintf(":%d", conf.Service.ServerPort))
	}

	// Listen for traffic indefinitely.
	if err := server.Serve(lis); err != nil {
		log.Fatalf("server ended: %s", err)
	}
}
