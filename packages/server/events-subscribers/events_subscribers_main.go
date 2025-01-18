package main

import (
	"context"
	"io"
	"log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/customeros/customeros/packages/server/events-subscribers/config"
	"github.com/customeros/customeros/packages/server/events-subscribers/handlers"
	"github.com/customeros/customeros/packages/server/events-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

const (
	AppName                                = "events-subscribers"
	EventsQueueName                        = "events"
	EventsFlowParticipantScheduleQueueName = "events-flow-participant-schedule"
)

func main() {
	ctx := context.Background()

	// Config
	cfg := config.Load()

	appLogger := logger.NewExtendedAppLogger(&cfg.Common.Infrastructure.LoggerConfig)
	appLogger.InitLogger()
	appLogger.WithName(AppName)

	// Initialize Tracing
	tracingCloser := initTracing(cfg, appLogger)
	if tracingCloser != nil {
		defer tracingCloser.Close()
	}
	defer tracing.RecoverAndLogToJaeger(appLogger)

	postgresDb, err := commonConfig.InitPostgres(&commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			PostgresConfig:      cfg.Common.Infrastructure.PostgresConfig,
			PostgresAsyncConfig: cfg.Common.Infrastructure.PostgresAsyncConfig,
		},
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	neo4jDriver, err := commonConfig.NewNeo4jDriver(cfg.Common.Infrastructure.Neo4jConfig)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Common.Infrastructure.Neo4jConfig.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Events processing
	var eventsProcessingGrpcClient *grpc_client.Clients
	if cfg.Common.Infrastructure.GrpcClientConfig.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(&cfg.Common.Infrastructure.GrpcClientConfig)
		gRPCconn, err := df.GetEventsProcessingPlatformConn()
		defer df.Close(gRPCconn)
		if err != nil {
			appLogger.Fatalf("Failed to connect: %v", err)
		}
		eventsProcessingGrpcClient = grpc_client.InitClients(gRPCconn)
	}

	postgresRepositories := postgres_repository.InitRepositories(postgresDb)
	neo4jRepositories := neo4j_repository.InitNeo4jRepositories(&neo4jDriver, cfg.Common.Infrastructure.Neo4jConfig.Database)

	commonServices := commonService.InitCommonServices(
		appLogger,
		neo4jRepositories,
		postgresRepositories,
		cfg.Common,
		eventsProcessingGrpcClient,
	)

	// Create dependencies for event handlers
	dependencies := &model.DependencyContainer{
		Logger:               appLogger,
		GRPCClients:          eventsProcessingGrpcClient,
		CommonConfig:         cfg.Common,
		PostgresRepositories: postgresRepositories,
		Neo4jRepositories:    neo4jRepositories,
		CommonServices:       commonServices,
	}

	// Create Events Service
	eventsService, err := events.NewEventsService(
		dependencies.CommonConfig.Infrastructure.RabbitMQConfig.Url,
		appLogger,
	)
	if err != nil {
		appLogger.Fatalf("Failed to create events service: %v", err)
	}
	defer eventsService.Close()

	// Register all handlers
	handlers.InitHandlerRegistration(eventsService, dependencies)

	// Set up queue listeners
	go func() {
		if err := eventsService.Subscriber.ListenQueue(EventsQueueName); err != nil {
			appLogger.Fatalf("Failed to listen to queue %s: %v", EventsQueueName, err)
		}
	}()

	go func() {
		if err := eventsService.Subscriber.ListenQueueExclusive(EventsFlowParticipantScheduleQueueName); err != nil {
			appLogger.Fatalf("Failed to listen to exclusive queue %s: %v", EventsFlowParticipantScheduleQueueName, err)
		}
	}()

	// Block the main thread from exiting
	forever := make(chan bool)
	log.Println(" [*] Waiting for messages")
	<-forever
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
