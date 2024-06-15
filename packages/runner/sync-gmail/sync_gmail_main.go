package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/caches"
	syncGmailConfig "github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	localCron "github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/cron"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	config := loadConfiguration()

	// Initialize logger
	appLogger := logger.NewExtendedAppLogger(&config.Logger)
	appLogger.InitLogger()
	appLogger.WithName("sync-gmail")

	// Tracing
	tracingCloser := initTracing(&config.Jaeger, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}

	sqlDb, gormDb, errPostgres := syncGmailConfig.NewPostgresClient(config)
	if errPostgres != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", errPostgres.Error())
	}
	defer sqlDb.Close()

	neo4jDriver, errNeo4j := syncGmailConfig.NewDriver(config)
	if errNeo4j != nil {
		logrus.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (*neo4jDriver).Close(ctx)

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(&config.GrpcClientConfig)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		logrus.Fatalf("failed opening connection to gRPC: %v", err.Error())
	}
	defer df.Close(gRPCconn)
	grpcContainer := grpc_client.InitClients(gRPCconn)

	appCache := caches.NewCache()
	services := service.InitServices(config, neo4jDriver, gormDb, grpcContainer, appCache)

	//init app cache
	personalEmailProviderEntities, err := services.Repositories.PostgresRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		appLogger.Fatalf("Error getting personal email providers: %s", err.Error())
	}
	personalEmailProviders := make([]string, 0)
	for _, personalEmailProvider := range personalEmailProviderEntities {
		personalEmailProviders = append(personalEmailProviders, personalEmailProvider.ProviderDomain)
	}
	appCache.SetPersonalEmailProviders(personalEmailProviders)

	emailExclusionEntities, err := services.Repositories.PostgresRepositories.EmailExclusionRepository.GetExclusionList()
	if err != nil {
		appLogger.Fatalf("Error getting email exclusion list: %s", err.Error())
	}
	appCache.SetEmailExclusion(emailExclusionEntities)

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

func loadConfiguration() *syncGmailConfig.Config {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Failed loading .env file")
	}

	cfg := syncGmailConfig.Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("%+v", err)
	}

	return &cfg
}

func initTracing(cfg *tracing.JaegerConfig, appLogger logger.Logger) io.Closer {
	if cfg.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(cfg, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		return closer
	}
	return nil
}
