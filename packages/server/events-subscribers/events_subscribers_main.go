package main

import (
	"context"
	"fmt"
	"io"
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
	"github.com/customeros/customeros/packages/server/events-subscribers/listeners"
	"github.com/customeros/customeros/packages/server/events-subscribers/logger"
	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

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
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:    ctx,
		cancel: cancel,
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

	// Initialize events service with configs
	publisherConfig := &events.PublisherConfig{
		MessageTTL:          events.DefaultMessageTTL,
		MaxRetries:          events.DefaultMaxRetries,
		PublishTimeout:      events.DefaultPublishTimeout,
		ReconnectBackoff:    events.DefaultReconnectBackoff,
		MaxReconnectBackoff: events.DefaultMaxReconnectBackoff,
	}

	subscriberConfig := &events.SubscriberConfig{
		MaxRetries:          events.DefaultMaxRetries,
		ReconnectBackoff:    events.DefaultReconnectBackoff,
		MaxReconnectBackoff: events.DefaultMaxReconnectBackoff,
	}

	eventsService, err := events.NewEventsService(
		a.config.Common.Infrastructure.RabbitMQConfig.Url,
		a.logger,
		publisherConfig,
		subscriberConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to create events service: %w", err)
	}
	a.events = eventsService

	// Register listeners instead of handlers
	if err := a.initializeListeners(); err != nil {
		return fmt.Errorf("failed to initialize listeners: %w", err)
	}

	return nil
}

func (a *App) initializeListeners() error {
	// Create and register listeners for each event type

	// Contact Listeners
	a.events.Subscriber.RegisterListener(listeners.NewAddSocialToContactListener(a.logger, a.deps))
	a.events.Subscriber.RegisterListener(listeners.NewHideContactListener(a.logger, a.deps))

	// Enrichment Listeners
	a.events.Subscriber.RegisterListener(listeners.NewRequestEnrichContactListener(a.logger, a.deps))
	a.events.Subscriber.RegisterListener(listeners.NewRequestValidateEmailListener(a.logger, a.deps))

	// SKU Listeners
	a.events.Subscriber.RegisterListener(listeners.NewSkuUpdateListener(a.logger, a.deps))

	return nil
}

func (a *App) Run() error {
	// Start listening to queues
	errChan := make(chan error, 4)

	// CustomerOS Events Queue
	go func() {
		if err := a.events.Subscriber.ListenQueue(events.QueueEvents); err != nil {
			errChan <- fmt.Errorf("failed to listen to events queue: %w", err)
		}
	}()

	// Agents Queue
	go func() {
		if err := a.events.Subscriber.ListenQueue(events.QueueAgents); err != nil {
			errChan <- fmt.Errorf("failed to listen to agents queue: %w", err)
		}
	}()

	// Notifications Queue
	go func() {
		if err := a.events.Subscriber.ListenQueue(events.QueueNotifications); err != nil {
			errChan <- fmt.Errorf("failed to listen to notifications queue: %w", err)
		}
	}()

	// Flow Participant Queue (exclusive)
	go func() {
		if err := a.events.Subscriber.ListenQueueExclusive(events.QueueFlowParticipantSchedule); err != nil {
			errChan <- fmt.Errorf("failed to listen to flow participant queue: %w", err)
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
