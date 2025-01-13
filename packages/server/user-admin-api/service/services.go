package service

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	commonConfig "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/azure"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/domain"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/email"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/google"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/industry"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_event"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/interaction_session"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/jobrole"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/opensrs"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/postmark"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/social"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/user"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
)

type Services struct {
	Log             logger.Logger
	Cache           *caches.Cache
	GrpcClients     *grpc_client.Clients
	Config          *config.Config
	Events          *events.EventsService
	Postgres        *repository.Repositories
	Neo4j           *neoRepo.Repositories
	MailService     interfaces.MailService
	PostmarkService interfaces.PostmarkService
}

func InitServices(cfg *config.Config, driver *neo4j.DriverWithContext, postgresDB *commonConfig.PostgresDB, grpcClients *grpc_client.Clients, cache *caches.Cache, appLogger logger.Logger) *Services {
	services := Services{
		Log:         appLogger,
		Cache:       cache,
		GrpcClients: grpcClients,
		Config:      cfg,
	}
	events, err := events.NewEventsService(cfg.RabbitMQ.Url, appLogger)
	if err != nil {
		log.Fatalf("cannot start rabbitMQ")
	}
	services.Events = events
	services.Postgres = repository.InitRepositories(postgresDB)
	services.Neo4j = neoRepo.InitNeo4jRepositories(driver, cfg.Neo4j.Database)

	interactionSession := interaction_session.NewInteractionSessionService(
		services.Neo4j,
	)

	postmark := postmark.NewPostmarkService(
		&cfg.Postmark,
		services.Postgres,
	)
	services.PostmarkService = postmark

	azure := azure.NewAzureService(
		&cfg.Azure,
		services.Postgres,
		services.Neo4j,
	)

	domain := domain.NewDomainService(
		services.Log,
		services.Cache,
		services.Postgres,
		services.Neo4j,
		services.Events,
	)

	google := google.NewGoogleService(
		&cfg.GoogleOAuth,
		services.Postgres,
		services.Neo4j,
	)

	opensrs := opensrs.NewOpenSRSService(
		services.Log,
		&services.Config.OpenSRS,
		services.Postgres,
	)

	industry := industry.NewIndustryService(
		services.Log,
		services.Neo4j,
	)

	user := user.NewUserService(
		services.Neo4j,
		services.Postgres,
		services.Events,
	)

	email := email.NewEmailService(
		services.Neo4j,
		services.Events,
		nil, // contact
		nil, // org
	)

	jobrole := jobrole.NewJobRoleService(
		services.Neo4j,
		services.Events,
		nil, // org
	)

	social := social.NewSocialService(
		services.Log,
		services.Neo4j,
		services.Events,
		nil, // contact
	)

	contact := contact.NewContactService(
		services.Log,
		services.Neo4j,
		services.Events,
		domain,
		nil, // email
		nil, // org
		nil, // jobrole
		nil, // social
	)

	org := organization.NewOrganizationService(
		services.Log,
		services.Postgres,
		services.Neo4j,
		services.Events,
		domain,
		industry,
		user,
		nil, // social
	)

	interactionEvent := interaction_event.NewInteractionEventService(
		services.Neo4j,
		nil, // email
	)

	mail := mail.NewMailService(
		cache,
		services.Postgres,
		services.Neo4j,
		azure,
		google,
		interactionSession,
		opensrs,
		nil, // contact
		nil, // email
		nil, // interactionEvent
		nil, // org
	)

	email.SetContactService(contact)
	email.SetOrganizationService(org)

	jobrole.SetOrganizationService(org)

	social.SetContactService(contact)

	contact.SetEmailService(email)
	contact.SetOrganizationService(org)
	contact.SetJobRoleService(jobrole)
	contact.SetSocialService(social)

	org.SetSocialService(social)

	interactionEvent.SetEmailService(email)

	mail.SetContactService(contact)
	mail.SetEmailService(email)
	mail.SetInteractionEventService(interactionEvent)
	mail.SetOrganizationService(org)

	// validate init
	if !email.IsInitialized() {
		logInitializationErr("Email")
	}

	if !jobrole.IsInitialized() {
		logInitializationErr("JobRole")
	}

	if !social.IsInitialized() {
		logInitializationErr("Social")
	}

	if !contact.IsInitialized() {
		logInitializationErr("Contact")
	}

	if !org.IsInitialized() {
		logInitializationErr("Organization")
	}

	if !interactionEvent.IsInitialized() {
		logInitializationErr("InteractionEvent")
	}

	if !mail.IsInitialized() {
		logInitializationErr("Mail")
	}

	return &services
}

func logInitializationErr(s string) {
	log.Fatalf(s, "service not fully initialized")
}
