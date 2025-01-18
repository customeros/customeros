package service

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
)

type Services struct {
	cfg *config.Config

	CommonServices     *commonService.CommonServices
	Neo4jRepository    *neo4jrepository.Repositories
	PostgresRepository *postgresRepository.Repositories

	TenantService             TenantService
	LocationService           LocationService
	UserService               UserService
	LogEntryService           LogEntryService
	OrganizationService       OrganizationService
	ContactService            ContactService
	IssueService              IssueService
	SyncStatusService         SyncStatusService
	ExternalSystemService     ExternalSystemService
	FinderService             FinderService
	InteractionSessionService InteractionSessionService
	CommentService            CommentService
	InvoiceService            InvoiceService
}

func InitServices(log logger.Logger,
	repositories *repository.Repositories,
	cfg *config.Config,
	commonServices *commonService.CommonServices,
	grpcClients *grpc_client.Clients,
	cache *caches.Cache,
) *Services {
	services := Services{
		CommonServices:     commonServices,
		Neo4jRepository:    repositories.Neo4jRepositories,
		PostgresRepository: repositories.PostgresRepositories,

		TenantService:             NewTenantService(log, repositories, cache),
		LocationService:           NewLocationService(log, repositories, grpcClients),
		SyncStatusService:         NewSyncStatusService(log, repositories),
		InteractionSessionService: NewInteractionSessionService(log, repositories, grpcClients),
	}
	services.cfg = cfg
	services.ExternalSystemService = NewExternalSystemService(log, repositories, cache, &services)
	services.UserService = NewUserService(log, repositories, grpcClients, &services)
	services.OrganizationService = NewOrganizationService(log, repositories, grpcClients, &services, cache)
	services.ContactService = NewContactService(log, repositories, grpcClients, &services)
	services.LogEntryService = NewLogEntryService(log, repositories, grpcClients, &services)
	services.IssueService = NewIssueService(log, repositories, grpcClients, &services)
	services.FinderService = NewFinderService(log, repositories, &services)
	services.CommentService = NewCommentService(log, repositories, grpcClients, &services)
	services.InvoiceService = NewInvoiceService(log, repositories, grpcClients, &services)
	return &services
}
