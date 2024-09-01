package service

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"gorm.io/gorm"
)

type Services struct {
	cfg *config.Config

	CommonServices *commonService.Services

	TenantService             TenantService
	EmailService              EmailService
	LocationService           LocationService
	PhoneNumberService        PhoneNumberService
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
	OrderService              OrderService
}

func InitServices(log logger.Logger,
	driver *neo4j.DriverWithContext,
	gormDB *gorm.DB,
	cfg *config.Config,
	commonServices *commonService.Services,
	grpcClients *grpc_client.Clients,
	cache *caches.Cache) *Services {
	repositories := repository.InitRepos(driver, gormDB, cfg.Neo4j.Database)

	services := Services{
		CommonServices:            commonServices,
		TenantService:             NewTenantService(log, repositories, cache),
		EmailService:              NewEmailService(log, repositories, grpcClients),
		LocationService:           NewLocationService(log, repositories, grpcClients),
		PhoneNumberService:        NewPhoneNumberService(log, repositories, grpcClients),
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
	services.OrderService = NewOrderService(log, repositories, grpcClients, &services)
	return &services
}
