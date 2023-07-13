package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	commsApiConfig "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/routes/ContactHub"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/service"
	"github.com/redis/go-redis/v9"
	"log"
)

func main() {
	config := loadConfiguration()

	graphqlClient := graphql.NewClient(config.Service.CustomerOsAPI)
	redisUrl := fmt.Sprintf("redis://%s", config.Redis.Host)
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("unvalid redis redisUrl: %s %v", redisUrl, err)
	}
	redisClient := redis.NewClient(opt)

	db, err := InitDB(&config)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	services := service.InitServices(graphqlClient, redisClient, &config, db)
	hub := ContactHub.NewContactHub()
	go hub.Run()
	routes.Run(&config, hub, services) // run this as a background goroutine

}

func InitDB(cfg *commsApiConfig.Config) (db *commsApiConfig.StorageDB, err error) {
	db, err = commsApiConfig.NewDBConn(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Db,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.MaxConn,
		cfg.Postgres.MaxIdleConn,
		cfg.Postgres.ConnMaxLifetime)
	if err != nil {
		return nil, fmt.Errorf("InitDB: Coud not open db connection: %s", err.Error())
	}
	return db, nil
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
