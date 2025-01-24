package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/events-subscribers/config"
	"github.com/customeros/customeros/packages/server/events-subscribers/handlers"
	"github.com/customeros/customeros/packages/server/events-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type QueueConfig struct {
	Events          string
	FlowParticipant string
}

type App struct {
	ctx           context.Context
	cancel        context.CancelFunc
	config        *config.Config
	logger        logger.Logger
	tracingCloser io.Closer
	deps          *model.DependencyContainer
	events        *events.EventsService
	postgresDB    *commonConfig.PostgresDB
	neo4jDriver   *neo4j.DriverWithContext
	queueConfig   QueueConfig
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
		queueConfig: QueueConfig{
			Events:          "events",
			FlowParticipant: "events-flow-participant-schedule",
		},
	}
}

func (a *App) Initialize() error {
	// Load config
	a.config = config.Load()

	// Initialize logger
	a.logger = logger.NewExtendedAppLogger(&a.config.Common.Infrastructure.LoggerConfig)
	a.logger.InitLogger()
	a.logger.WithName("events-subscribers")

	// Initialize tracing
	if err := a.initTracing(); err != nil {
		return fmt.Errorf("failed to initialize tracing: %w", err)
	}

	// Initialize databases
	if err := a.initDatabases(); err != nil {
		return fmt.Errorf("failed to initialize databases: %w", err)
	}

	// Initialize services
	if err := a.initServices(); err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	return nil
}

func (a *App) initTracing() error {
	if a.config.Common.Infrastructure.JaegerConfig.Enabled {
		tracer, closer, err := tracing.NewJaegerTracer(&a.config.Common.Infrastructure.JaegerConfig, a.logger)
		if err != nil {
			return fmt.Errorf("could not initialize jaeger tracer: %w", err)
		}
		opentracing.SetGlobalTracer(tracer)
		a.tracingCloser = closer
	}
	return nil
}

func (a *App) initDatabases() error {
	// Initialize Postgres
	db, err := commonConfig.InitPostgres(&commonConfig.CommonConfig{
		Infrastructure: commonConfig.InfrastructureConfig{
			PostgresConfig:      a.config.Common.Infrastructure.PostgresConfig,
			PostgresAsyncConfig: a.config.Common.Infrastructure.PostgresAsyncConfig,
		},
	})
	if err != nil {
		return fmt.Errorf("failed opening connection to postgres: %w", err)
	}
	a.postgresDB = db

	// Initialize Neo4j
	driver, err := commonConfig.NewNeo4jDriver(a.config.Common.Infrastructure.Neo4jConfig)
	if err != nil {
		return fmt.Errorf("could not establish connection with neo4j: %w", err)
	}
	a.neo4jDriver = &driver

	return nil
}

func (a *App) initServices() error {
	// Initialize gRPC clients
	var eventsProcessingGrpcClient *grpc_client.Clients
	if a.config.Common.Infrastructure.GrpcClientConfig.EventsProcessingPlatformEnabled {
		df := grpc_client.NewDialFactory(&a.config.Common.Infrastructure.GrpcClientConfig)
		gRPCconn, err := df.GetEventsProcessingPlatformConn()
		if err != nil {
			return fmt.Errorf("failed to connect to gRPC: %w", err)
		}
		eventsProcessingGrpcClient = grpc_client.InitClients(gRPCconn)
	}

	// Initialize repositories
	postgresRepositories := postgres_repository.InitRepositories(a.postgresDB)
	neo4jRepositories := neo4j_repository.InitNeo4jRepositories(a.neo4jDriver, a.config.Common.Infrastructure.Neo4jConfig.Database)

	// Initialize common services
	commonServices := service.InitCommonServices(
		a.logger,
		neo4jRepositories,
		postgresRepositories,
		a.config.Common,
		eventsProcessingGrpcClient,
		&service.InitOptions{LoadPersonalEmailProviders: true},
	)

	// Initialize dependency container
	a.deps = &model.DependencyContainer{
		Logger:               a.logger,
		GRPCClients:          eventsProcessingGrpcClient,
		CommonConfig:         a.config.Common,
		PostgresRepositories: postgresRepositories,
		Neo4jRepositories:    neo4jRepositories,
		CommonServices:       commonServices,
	}

	// Initialize events service
	eventsService, err := events.NewEventsService(
		a.config.Common.Infrastructure.RabbitMQConfig.Url,
		a.logger,
	)
	if err != nil {
		return fmt.Errorf("failed to create events service: %w", err)
	}
	a.events = eventsService

	// Initialize handlers
	handlers.InitHandlerRegistration(a.events, a.deps)

	return nil
}

func (a *App) Run() error {
	// Start event listeners
	errChan := make(chan error, 2)

	go func() {
		if err := a.events.Subscriber.ListenQueue(a.queueConfig.Events); err != nil {
			errChan <- fmt.Errorf("failed to listen to queue %s: %w", a.queueConfig.Events, err)
		}
	}()

	go func() {
		if err := a.events.Subscriber.ListenQueueExclusive(a.queueConfig.FlowParticipant); err != nil {
			errChan <- fmt.Errorf("failed to listen to exclusive queue %s: %w", a.queueConfig.FlowParticipant, err)
		}
	}()

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		return a.Shutdown()
	case <-a.ctx.Done():
		return a.Shutdown()
	}
}

func (a *App) Shutdown() error {
	a.cancel()

	if a.tracingCloser != nil {
		a.tracingCloser.Close()
	}
	if a.postgresDB != nil {
		a.postgresDB.Close()
	}
	if a.neo4jDriver != nil {
		(*a.neo4jDriver).Close(a.ctx)
	}
	if a.events != nil {
		a.events.Close()
	}

	return nil
}

func main() {
	app := NewApp()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

