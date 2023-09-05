package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	syncGmailRawConfig "github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	localCron "github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/cron"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/service"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	config := loadConfiguration()

	sqlDb, gormDb, errPostgres := syncGmailRawConfig.NewPostgresClient(config)
	if errPostgres != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	neo4jDriver, errNeo4j := syncGmailRawConfig.NewDriver(config)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	services := service.InitServices(neo4jDriver, gormDb, config)

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName("sync-gmail-raw")

	cronJub := localCron.StartCron(config, services)

	if err := run(appLogger, cronJub); err != nil {
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
