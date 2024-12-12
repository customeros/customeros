package main

import (
	"context"
	"io"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/listeners"
	"github.com/openline-ai/openline-customer-os/packages/server/events-subscribers/logger"
)

const (
	AppName = "events-subscribers"
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

	// Register listeners
	commonServices.RabbitMQService.RegisterHandler(dto.FlowOn{}, listeners.Handle_FlowOn)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowParticipantSchedule{}, listeners.Handle_FlowParticipantSchedule)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowComputeParticipantsRequirements{}, listeners.Handle_FlowComputeParticipantsRequirements)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowParticipantGoalAchieved{}, listeners.Handle_FlowParticipantGoalAchieved)

	// mailstack
	commonServices.RabbitMQService.RegisterHandler(dto.MailstackProvisionBuyRequest{}, listeners.Handle_MailstackProvisionBuyRequest)
	commonServices.RabbitMQService.RegisterHandler(dto.MailstackProvisionMailbox{}, listeners.Handle_MailstackProvisionMailbox)

	// contact
	commonServices.RabbitMQService.RegisterHandler(dto.AddSocialToContact{}, listeners.OnSocialAddedToContact)
	commonServices.RabbitMQService.RegisterHandler(dto.RequestEnrichContact{}, listeners.OnRequestedEnrichContact)

	// organization
	commonServices.RabbitMQService.RegisterHandler(dto.RequestRefreshLastTouchpoint{}, listeners.OnRequestLastTouchpointRefresh)
	commonServices.RabbitMQService.RegisterHandler(dto.RequestEnrichOrganization{}, listeners.OnRequestedEnrichOrganization)

	// email
	commonServices.RabbitMQService.RegisterHandler(dto.RequestValidateEmail{}, listeners.OnRequestedValidateEmail)

	// FlowEngine
	commonServices.RabbitMQService.RegisterHandler(dto.WebhookEvent{}, listeners.OnWebhookEventCreated)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowActionEvent{}, listeners.OnFlowActionEventCreated)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowActionExecutionResultEvent{}, listeners.OnFlowActionExecutionResultsEventCreated)

	// Listen for messages
	commonServices.RabbitMQService.ListenQueue(commonService.EventsQueueName)
	commonServices.RabbitMQService.ListenQueueExclusive(commonService.EventsFlowParticipantScheduleQueueName)

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
