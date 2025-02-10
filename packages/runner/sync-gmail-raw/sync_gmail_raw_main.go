package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	syncGmailRawConfig "github.com/customeros/customeros/packages/runner/sync-gmail-raw/config"
	localCron "github.com/customeros/customeros/packages/runner/sync-gmail-raw/cron"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/logger"
	"github.com/customeros/customeros/packages/runner/sync-gmail-raw/service"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
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

	neo4jDriver, errNeo4j := syncGmailRawConfig.NewDriver(config)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName("sync-gmail-raw")

	services := service.InitServices(neo4jDriver, postgresDb, config, appLogger)

	all, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetAll(ctx)

	for _, v := range all {

		v.AccessToken, err = postgres_entity.EncryptToken(config.GoogleOAuthConfig.EncryptionKey, v.AccessToken)
		if err != nil {
			appLogger.Error(err)
		}
		v.RefreshToken, err = postgres_entity.EncryptToken(config.GoogleOAuthConfig.EncryptionKey, v.RefreshToken)
		if err != nil {
			appLogger.Error(err)
		}
		_, err := services.CommonServices.PostgresRepositories.OAuthTokenRepository.Update(ctx, v.TenantName, v.PlayerIdentityId, v.Provider, v.AccessToken, v.RefreshToken, v.ExpiresAt)
		if err != nil {
			appLogger.Error(err)
		}

	}

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
