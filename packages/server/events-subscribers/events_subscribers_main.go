package main

import (
	"context"
	"io"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
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

	// Register listeners
	commonServices.RabbitMQService.RegisterHandler(dto.Delete{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateBankAccount{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.DeleteBankAccount{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateBankAccount{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateComment{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateComment{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateContact{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateContact{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddContactToOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateSocialForContact{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveSocialFromContact{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateContract{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.ChangeStatusForContract{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateContract{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateCustomFieldTemplate{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateCustomFieldTemplate{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateDomain{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddDomain{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveDomain{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.AddEmail{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RegisterEmail{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveEmail{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateOrganizationOnboardingStatus{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddActionToOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddParentOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddSocialToOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddSubOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.MergeOrganizations{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveParentOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveSocialFromOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveSubOrganization{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateSocialForOrganization{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.AddUserAssigneeToIssue{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.AddUserFollowerToIssue{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.CreateIssue{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveUserAssigneeFromIssue{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveUserFollowerFromIssue{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateIssue{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.SaveJobRole{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateLocation{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateLogEntry{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateLogEntry{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateOpportunity{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateOpportunity{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CloseServiceLineItem{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.CreateServiceLineItem{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.DeleteServiceLineItem{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.PauseServiceLineItem{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.ResumeServiceLineItem{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateServiceLineItem{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateSocial{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateSocial{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.AddTag{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.RemoveTag{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.SaveTag{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateTenantBillingProfile{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateTenantBillingProfile{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateTenantSettings{}, listeners.Handle_NotHandledListener)

	commonServices.RabbitMQService.RegisterHandler(dto.CreateUser{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UserLogin{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.UpdateUser{}, listeners.Handle_NotHandledListener)

	// flow
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
	commonServices.RabbitMQService.RegisterHandler(dto.ShowContact{}, listeners.Handle_NotHandledListener)
	commonServices.RabbitMQService.RegisterHandler(dto.HideContact{}, listeners.OnContactHidden)

	// organization
	commonServices.RabbitMQService.RegisterHandler(dto.RequestRefreshLastTouchpoint{}, listeners.OnRequestLastTouchpointRefresh)
	commonServices.RabbitMQService.RegisterHandler(dto.RequestEnrichOrganization{}, listeners.OnRequestedEnrichOrganization)

	// email
	commonServices.RabbitMQService.RegisterHandler(dto.RequestValidateEmail{}, listeners.OnRequestedValidateEmail)

	// Automation Engine
	commonServices.RabbitMQService.RegisterHandler(dto.WebhookEvent{}, listeners.OnWebhookEventCreated)

	// Flow Engine
	commonServices.RabbitMQService.RegisterHandler(dto.FlowAgentEvent{}, listeners.OnFlowAgentEventCreated)
	commonServices.RabbitMQService.RegisterHandler(dto.FlowAgentExecutionResultEvent{}, listeners.OnFlowAgentExecutionResultsEventCreated)

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
