package main

import (
	"context"
	"io"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/model"
)

const (
	AppName                                = "events-subscribers"
	EventsQueueName                        = "events"
	EventsFlowParticipantScheduleQueueName = "events-flow-participant-schedule"
)

func main() {
	ctx := context.Background()

	cfg := loadConfiguration()

	appLogger := logger.NewExtendedAppLogger(&cfg.Logger)
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
			PostgresConfig:      cfg.PostgresConfig,
			PostgresAsyncConfig: cfg.PostgresAsyncConfig,
		},
	})
	if err != nil {
		logrus.Fatalf("failed opening connection to postgres: %v", err.Error())
	}
	defer postgresDb.Close()

	neo4jDriver, err := commonConfig.NewNeo4jDriver(cfg.Neo4j)
	if err != nil {
		appLogger.Fatalf("Could not establish connection with neo4j at: %v, error: %v", cfg.Neo4j.Target, err.Error())
	}
	defer neo4jDriver.Close(ctx)

	// Events processing
	var eventsProcessingGrpcClient *grpc_client.Clients
	if cfg.GrpcClientConfig.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(&cfg.GrpcClientConfig)
		gRPCconn, err := df.GetEventsProcessingPlatformConn()
		defer df.Close(gRPCconn)
		if err != nil {
			appLogger.Fatalf("Failed to connect: %v", err)
		}
		eventsProcessingGrpcClient = grpc_client.InitClients(gRPCconn)
	}

	postgresRepositories := repository.InitRepositories(postgresDb)
	neo4jRepositories := neo4jRepo.InitNeo4jRepositories(&neo4jDriver, cfg.Neo4j.Database)

	commonServices := commonService.InitCommonServices(
		appLogger,
		neo4jRepositories,
		postgresRepositories,
		&cfg.CommonConfig,
		eventsProcessingGrpcClient,
	)

	// Create dependencies for event handlers
	dependencies := &model.DependencyContainer{
		Logger:               appLogger,
		GRPCClients:          eventsProcessingGrpcClient,
		CommonConfig:         &cfg.CommonConfig,
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

func loadConfiguration() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}

func initTracing(cfg *config.Config, appLogger logger.Logger) io.Closer {
	if cfg.Jaeger.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&cfg.Jaeger, appLogger)
		if err != nil {
			appLogger.Fatalf("Could not initialize jaeger tracer: %v", err.Error())
		}
		opentracing.SetGlobalTracer(tracer)
		return closer
	}
	return nil
}
