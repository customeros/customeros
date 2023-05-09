package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/config"
	msProto "github.com/openline-ai/openline-customer-os/packages/server/message-store-api/proto/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/service"
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

	neo4jDriver, err := config.NewDriver(conf)
	if err != nil {
		log.Fatalf("failed opening connection to neo4j: %v", err.Error())
	}
	ctx := context.Background()
	defer (*neo4jDriver).Close(ctx)

	graphqlClient := graphql.NewClient(conf.Service.CustomerOsAPI)

	// Create a new gRPC server (you can wire multiple services to a single server).
	server := grpc.NewServer()

	repositories := repository.InitRepositories(db.GormDB, neo4jDriver)
	commonStoreService := service.NewCommonStoreService()
	customerOSService := service.NewCustomerOSService(neo4jDriver, graphqlClient, repositories, commonStoreService, conf)

	// Register the Message Item service with the server.
	msProto.RegisterMessageStoreServiceServer(server, service.NewMessageService(neo4jDriver, repositories, customerOSService, commonStoreService))

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
