package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"

	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	agent_producers "github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_event_producers"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/temporal/worker"
	"github.com/opentracing/opentracing-go"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/config"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/container"
	localcron "github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/cron"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/repository"
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
	defer telemetry.RecoverAndLogMain(appLogger)

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

	repositories := repository.InitRepositories(cfg, &neo4jDriver, postgresDb)

	// Check if migration is requested
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		appLogger.Info("Running database migration...")
		if err := repositories.PostgresRepositories.MigrateOpenlineDB(); err != nil {
			appLogger.Fatalf("Database migration failed: %v", err)
		}
		appLogger.Info("Database migration completed successfully")
	}

	cntnr := &container.Container{
		Cfg:            cfg,
		Log:            appLogger,
		Repositories:   repositories,
		CommonServices: commonService.InitCommonServices(appLogger, repositories.Neo4jRepositories, repositories.PostgresRepositories, cfg.Common, nil, &commonService.InitOptions{LoadPersonalEmailProviders: true}),
	}
	cntnr.AgentProducers = agent_producers.InitAgentProducers(cntnr.CommonServices)

	// Initialize waitGroup for goroutines
	var waitGroup sync.WaitGroup

	// Run Temporal worker
	if cfg.Common.External.TemporalConfig.RunWorker {
		waitGroup.Add(1)
		go runTemporalWorker(cfg, appLogger, &waitGroup)
	}

	crons := localcron.StartCron(cntnr)

	if err = run(appLogger, crons,
		func() error {
			return localcron.StopCron(appLogger, crons) // Stop cron jobs
		}); err != nil {
		appLogger.Fatal(err)
	}

	// Wait for all goroutines to complete
	waitGroup.Wait()

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
	var closer io.Closer

	// Initialize Jaeger if enabled
	if cfg.Common.Infrastructure.JaegerConfig.Enabled {
		tracer, jaegerCloser, err := telemetry.NewJaegerTracer(&cfg.Common.Infrastructure.JaegerConfig, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		closer = jaegerCloser
	}

	// Initialize OpenTelemetry
	err := telemetry.InitOpenTelemetry(context.Background(), &cfg.Common.Infrastructure.OpenTelemetryConfig)
	if err != nil {
		appLogger.Warnf("Could not initialize OpenTelemetry: %v", err.Error())
	}

	return closer
}

func runTemporalWorker(cfg *config.Config, logger logger.Logger, waitGroup *sync.WaitGroup) {
	// Start it in the background
	go func() {
		if err := worker.RunWebhookWorker(cfg.Common.External.TemporalConfig.HostPort, cfg.Common.External.TemporalConfig.Namespace); err != nil {
			logger.Error(err)
		}
		waitGroup.Done()
	}()
}
