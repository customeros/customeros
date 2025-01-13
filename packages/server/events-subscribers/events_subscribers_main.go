package main

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/listeners"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/logger"
)

const (
	AppName                                = "events-subscribers"
	EventsQueueName                        = "events"
	EventsFlowParticipantScheduleQueueName = "events-flow-participant-schedule"
)

// HandleDependencies contains all dependencies that might be needed by event handlers
type HandleDependencies struct {
	Logger         logger.Logger
	PostgresDB     *sql.DB
	Neo4jDriver    neo4j.DriverWithContext
	CommonServices *commonService.Services
	GRPCClients    *grpc_client.Clients
}

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

	postgresDb, err := commonConfig.InitPostgres(&commonConfig.GlobalConfig{
		PostgresConfig:      &cfg.PostgresConfig,
		PostgresAsyncConfig: &cfg.PostgresAsyncConfig,
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

	commonServices := commonService.InitServices(&commonConfig.GlobalConfig{
		RabbitMQConfig: &cfg.RabbitMQ,
		NovuConfig:     &cfg.NovuConfig,
		InternalServices: commonConfig.InternalServices{
			EnrichmentApiConfig: cfg.InternalServices.EnrichmentApi,
			AiApiConfig:         cfg.InternalServices.AiApi,
			ValidationApiConfig: cfg.InternalServices.ValidationApi,
		},
		ExternalServices: commonConfig.ExternalServices{
			OpenSRSConfig:    cfg.OpenSRSConfig,
			NamecheapConfig:  cfg.NamecheapConfig,
			CloudflareConfig: cfg.CloudflareConfig,
		},
	}, postgresDb, &neo4jDriver, cfg.Neo4j.Database, eventsProcessingGrpcClient, appLogger)

	// Create dependencies for event handlers
	dependencies := &HandleDependencies{
		Logger:         appLogger,
		PostgresDB:     postgresDb,
		Neo4jDriver:    *neo4jDriver,
		CommonServices: commonServices,
		GRPCClients:    eventsProcessingGrpcClient,
	}

	// Create Events Service
	eventsService, err := events.NewEventsService(
		events.Config{
			URL: cfg.RabbitMQ.Url,
		},
		appLogger,
	)
	if err != nil {
		appLogger.Fatalf("Failed to create events service: %v", err)
	}
	defer eventsService.Close()

	// Register Flow handlers
	eventsService.RegisterHandler(dto.FlowOn{},
		events.NewHandler(dto.FlowOn{},
			func(ctx context.Context, event dto.FlowOn) error {
				return listeners.Handle_FlowOn(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.FlowParticipantSchedule{},
		events.NewHandler(dto.FlowParticipantSchedule{},
			func(ctx context.Context, event dto.FlowParticipantSchedule) error {
				return listeners.Handle_FlowParticipantSchedule(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.FlowComputeParticipantsRequirements{},
		events.NewHandler(dto.FlowComputeParticipantsRequirements{},
			func(ctx context.Context, event dto.FlowComputeParticipantsRequirements) error {
				return listeners.Handle_FlowComputeParticipantsRequirements(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.FlowParticipantGoalAchieved{},
		events.NewHandler(dto.FlowParticipantGoalAchieved{},
			func(ctx context.Context, event dto.FlowParticipantGoalAchieved) error {
				return listeners.Handle_FlowParticipantGoalAchieved(ctx, dependencies, event)
			},
		),
	)

	// Mailstack handlers
	eventsService.RegisterHandler(dto.MailstackProvisionBuyRequest{},
		events.NewHandler(dto.MailstackProvisionBuyRequest{},
			func(ctx context.Context, event dto.MailstackProvisionBuyRequest) error {
				return listeners.Handle_MailstackProvisionBuyRequest(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.MailstackProvisionMailbox{},
		events.NewHandler(dto.MailstackProvisionMailbox{},
			func(ctx context.Context, event dto.MailstackProvisionMailbox) error {
				return listeners.Handle_MailstackProvisionMailbox(ctx, dependencies, event)
			},
		),
	)

	// Contact handlers
	eventsService.RegisterHandler(dto.AddSocialToContact{},
		events.NewHandler(dto.AddSocialToContact{},
			func(ctx context.Context, event dto.AddSocialToContact) error {
				return listeners.OnSocialAddedToContact(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.RequestEnrichContact{},
		events.NewHandler(dto.RequestEnrichContact{},
			func(ctx context.Context, event dto.RequestEnrichContact) error {
				return listeners.OnRequestedEnrichContact(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.HideContact{},
		events.NewHandler(dto.HideContact{},
			func(ctx context.Context, event dto.HideContact) error {
				return listeners.OnContactHidden(ctx, dependencies, event)
			},
		),
	)

	// Organization handlers
	eventsService.RegisterHandler(dto.RequestRefreshLastTouchpoint{},
		events.NewHandler(dto.RequestRefreshLastTouchpoint{},
			func(ctx context.Context, event dto.RequestRefreshLastTouchpoint) error {
				return listeners.OnRequestLastTouchpointRefresh(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.RequestEnrichOrganization{},
		events.NewHandler(dto.RequestEnrichOrganization{},
			func(ctx context.Context, event dto.RequestEnrichOrganization) error {
				return listeners.OnRequestedEnrichOrganization(ctx, dependencies, event)
			},
		),
	)

	// Email handlers
	eventsService.RegisterHandler(dto.RequestValidateEmail{},
		events.NewHandler(dto.RequestValidateEmail{},
			func(ctx context.Context, event dto.RequestValidateEmail) error {
				return listeners.OnRequestedValidateEmail(ctx, dependencies, event)
			},
		),
	)

	// Flow Engine handlers
	eventsService.RegisterHandler(dto.WebhookEvent{},
		events.NewHandler(dto.WebhookEvent{},
			func(ctx context.Context, event dto.WebhookEvent) error {
				return listeners.OnWebhookEventCreated(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.FlowAgentEvent{},
		events.NewHandler(dto.FlowAgentEvent{},
			func(ctx context.Context, event dto.FlowAgentEvent) error {
				return listeners.OnFlowAgentEventCreated(ctx, dependencies, event)
			},
		),
	)
	eventsService.RegisterHandler(dto.FlowAgentExecutionResultEvent{},
		events.NewHandler(dto.FlowAgentExecutionResultEvent{},
			func(ctx context.Context, event dto.FlowAgentExecutionResultEvent) error {
				return listeners.OnFlowAgentExecutionResultsEventCreated(ctx, dependencies, event)
			},
		),
	)

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
