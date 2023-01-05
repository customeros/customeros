package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/config"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/gen"
	pb "github.com/openline-ai/openline-customer-os/packages/server/message-store/gen/proto"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func InitDB(cfg *config.Config) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Db,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.MaxConn,
		cfg.Postgres.MaxIdleConn,
		cfg.Postgres.ConnMaxLifetime); err != nil {
		log.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func main() {
	conf := loadConfiguration()

	//GORM
	db, _ := InitDB(conf)
	defer db.SqlDB.Close()

	//ENT
	var connUrl = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", conf.Postgres.Host, conf.Postgres.Port, conf.Postgres.User, conf.Postgres.Db, conf.Postgres.Password)
	log.Printf("Connecting to database %s", connUrl)
	client, err := gen.Open("postgres", connUrl)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	neo4jDriver, err := config.NewDriver(conf)
	if err != nil {
		log.Fatalf("failed opening connection to neo4j: %v", err.Error())
	}
	defer (*neo4jDriver).Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	graphqlClient := graphql.NewClient(conf.Service.CustomerOsAPI)

	// Initialize the generated User service.
	svc := service.NewMessageService(client, neo4jDriver, graphqlClient, repository.InitRepositories(db.GormDB), conf)

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

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v", err)
	}

	return &cfg
}
