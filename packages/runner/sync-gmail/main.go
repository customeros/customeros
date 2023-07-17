package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/service"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	cfg := loadConfiguration()
	config.InitLogger(cfg)

	sqlDb, gormDb, errPostgres := config.NewPostgresClient(cfg)
	if errPostgres != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	neo4jDriver, errNeo4j := config.NewDriver(cfg)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	services := service.InitServices(neo4jDriver, gormDb)

	tenants, err := services.TenantService.GetAllTenants(ctx)
	if err != nil {
		panic(err)
	}

	for _, tenant := range tenants {

		usersForTenant, err := services.UserService.GetAllUsersForTenant(ctx, tenant.Name)
		if err != nil {
			panic(err)
		}

		for _, user := range usersForTenant {
			logrus.Infof("user: %v", user)
		}

	}

	services.EmailService.ReadNewEmailsForUsername("openline", "edi@openline.ai")

	//job - read all users and trigger email sync per user ( 5 mintues )
	//job - read all new emails for a user and sync them ( 1 minute )
	//1 job per user in thread pool

	//job 1
	//get tenants
	//get users in tenant
	//sync emails

	//get all users from tenant with access enabled to gmail ( all users except blacklisted )
}

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Failed loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v", err)
	}

	return &cfg
}
