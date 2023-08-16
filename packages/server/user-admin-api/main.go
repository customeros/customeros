package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	authCommonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/service"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/routes"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
	"log"
)

const (
	AppName = "USER-ADMIN-API"
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

func InitDB(cfg *config.Config, appLogger logger.Logger) (db *config.StorageDB, err error) {
	if db, err = config.NewDBConn(cfg); err != nil {
		appLogger.Fatalf("Coud not open db connection: %s", err.Error())
	}
	return
}

func main() {
	config := loadConfiguration()

	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName(AppName)

	db, _ := InitDB(&config, appLogger)
	defer db.SqlDB.Close()

	authServices := authCommonService.InitServices(db.GormDB)

	graphqlClient := graphql.NewClient(config.CustomerOS.CustomerOsAPI)
	cosClient := service.NewCustomerOsClient(&config, graphqlClient)

	routes.Run(&config, cosClient, authServices)
}
