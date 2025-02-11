package main

import (
	"context"
	"github.com/caarlos0/env/v6"
	syncGmailConfig "github.com/customeros/customeros/packages/runner/sync-gmail/config"
	localCron "github.com/customeros/customeros/packages/runner/sync-gmail/cron"
	"github.com/customeros/customeros/packages/runner/sync-gmail/logger"
	"github.com/customeros/customeros/packages/runner/sync-gmail/service"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/joho/godotenv"
	"github.com/opentracing/opentracing-go"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io"
	"log"
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
	defer tracing.RecoverAndLogToJaeger(appLogger)

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

	// Setting up gRPC client
	df := grpc_client.NewDialFactory(&config.GrpcClientConfig)
	gRPCconn, err := df.GetEventsProcessingPlatformConn()
	if err != nil {
		logrus.Fatalf("failed opening connection to gRPC: %v", err.Error())
	}
	defer df.Close(gRPCconn)

	services := service.InitServices(config, &neo4jDriver, postgresDb, appLogger)

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
