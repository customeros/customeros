package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	syncGmailRawConfig "github.com/customeros/customeros/packages/runner/sync-gmail-raw/config"
	localCron "github.com/customeros/customeros/packages/runner/sync-gmail-raw/cron"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/logger"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/service"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	config := loadConfiguration()

	postgresDb, err := commonConfig.InitPostgres(&commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			PostgresConfig:      config.PostgresConfig,
			PostgresAsyncConfig: config.PostgresAsyncConfig,
		},
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	neo4jDriver, err := commonConfig.NewNeo4jDriver(config.Neo4jDb)
	if err != nil {
		log.Fatalf("Could not establish connection with neo4j at: %v, error: %v", config.Neo4jDb.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName("sync-gmail-raw")

	services := service.InitServices(&neo4jDriver, postgresDb, config, appLogger)

	cronJobs := localCron.StartCronJobs(config, services)

	if err := run(appLogger, cronJobs); err != nil {
		appLogger.Fatal(err)
	}

	// Flush logs and exit
	appLogger.Sync()
}

func run(log logger.Logger, cron *cron.Cron) error {
	defer cron.Stop()

	// Shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdown
	log.Infof("Received shutdown signal %v", sig)

	// Gracefully stop
	if err := localCron.StopCron(log, cron); err != nil {
		return err
	}
	log.Info("Graceful shutdown complete")

	return nil
}

func loadConfiguration() *syncGmailRawConfig.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Failed loading .env file")
	}

	cfg := syncGmailRawConfig.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v", err)
	}

	return &cfg
}
