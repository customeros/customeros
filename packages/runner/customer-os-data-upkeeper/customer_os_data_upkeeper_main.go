package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/container"
	localcron "github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/cron"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/events/eventbuffer"
	"github.com/opentracing/opentracing-go"
)

func main() {
	// Config
	cfg := config.Load()

	// Logging
	appLogger := initLogger(cfg)

	// Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}
	defer tracing.RecoverAndLogToJaeger(appLogger)

	ctx := context.Background()

	// Initialize postgres db
	postgresDb, err := commonConfig.InitPostgres(cfg.Common)
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	// Neo4j DB
	neo4jDriver, errNeo4j := commonConfig.NewNeo4jDriver(cfg.Common.Infrastructure.Neo4jConfig)
	if errNeo4j != nil {
		appLogger.Fatalf("failed opening connection to neo4j: %v", errNeo4j.Error())
	}
	defer (neo4jDriver).Close(ctx)

	// Events processing
	var epClient *grpc_client.Clients
	if cfg.Common.Infrastructure.GrpcClientConfig.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(&cfg.Common.Infrastructure.GrpcClientConfig)
		gRPCconn, err := df.GetEventsProcessingPlatformConn()
		defer df.Close(gRPCconn)
		if err != nil {
			appLogger.Fatalf("Failed to connect: %v", err)
		}
		epClient = grpc_client.InitClients(gRPCconn)
	}

	repositories := repository.InitRepositories(cfg, &neo4jDriver, postgresDb)

	eventBufferProcessService := eventbuffer.NewEventBufferProcessService(repositories.PostgresRepositories.EventBufferRepository, appLogger, epClient)
	eventBufferProcessService.Start(ctx)
	defer eventBufferProcessService.Stop()

	eventBufferStoreService := eventbuffer.NewEventBufferStoreService(repositories.PostgresRepositories.EventBufferRepository, appLogger)

	cntnr := &container.Container{
		Cfg:                           cfg,
		Log:                           appLogger,
		Repositories:                  repositories,
		CommonServices:                commonService.InitCommonServices(appLogger, repositories.Neo4jRepositories, repositories.PostgresRepositories, cfg.Common, epClient, &commonService.InitOptions{LoadPersonalEmailProviders: true}),
		EventProcessingServicesClient: epClient,
		EventBufferStoreService:       eventBufferStoreService,
	}

	crons := localcron.StartCron(cntnr)

	if err = run(appLogger, crons,
		func() error {
			return localcron.StopCron(appLogger, crons) // Stop cron jobs
		},
		func() error {
			eventBufferProcessService.Stop() // Stop event buffer service
			return nil
		}); err != nil {
		appLogger.Fatal(err)
	}

	// Flush logs and exit
	appLogger.Sync()
}

func run(log logger.Logger, cron *cron.Cron, cleanupTasks ...func() error) error {
	defer cron.Stop() // Stop cron jobs first

	// Shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdown
	log.Infof("Received shutdown signal %v", sig)

	// Run cleanup tasks
	for _, task := range cleanupTasks {
		if err := task(); err != nil {
			log.Errorf("Cleanup task failed: %v", err)
		}
	}

	log.Info("Graceful shutdown complete")
	return nil
}

func initLogger(cfg *config.Config) logger.Logger {
	appLogger := logger.NewExtendedAppLogger(&cfg.Common.Infrastructure.LoggerConfig)
	appLogger.InitLogger()
	appLogger.WithName(constants.ServiceName)
	return appLogger
}

func initTracing(cfg *config.Config, appLogger logger.Logger) io.Closer {
	if cfg.Common.Infrastructure.JaegerConfig.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&cfg.Common.Infrastructure.JaegerConfig, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		return closer
	}
	return nil
}
