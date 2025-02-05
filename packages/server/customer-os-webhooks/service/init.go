package service

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonService "github.com/customeros/customeros/packages/server/customer-os-common-module/services"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-webhooks/caches"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/config"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
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
	cache *caches.Cache,
) *Services {
	services := Services{
		CommonServices:     commonServices,
		Neo4jRepository:    repositories.Neo4jRepositories,
		PostgresRepository: repositories.PostgresRepositories,

		TenantService:             NewTenantService(log, repositories, cache),
		LocationService:           NewLocationService(log, repositories),
		SyncStatusService:         NewSyncStatusService(log, repositories),
		InteractionSessionService: NewInteractionSessionService(log, repositories),
	}
	services.cfg = cfg
	services.ExternalSystemService = NewExternalSystemService(log, repositories, cache, &services)
	services.UserService = NewUserService(log, repositories, &services)
	services.OrganizationService = NewOrganizationService(log, repositories, &services, cache)
	services.ContactService = NewContactService(log, repositories, &services)
	services.LogEntryService = NewLogEntryService(log, repositories, &services)
	services.IssueService = NewIssueService(log, repositories, &services)
	services.FinderService = NewFinderService(log, repositories, &services)
	services.CommentService = NewCommentService(log, repositories, &services)
	services.InvoiceService = NewInvoiceService(log, repositories, &services)
	return &services
}
