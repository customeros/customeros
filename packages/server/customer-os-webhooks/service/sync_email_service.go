package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
)

type SyncEmailService interface {
	SyncEmail(ctx context.Context, email model.EmailData) (SyncResult, error)
}

type syncEmailService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
	maxWorkers   int
}

func NewSyncEmailService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) SyncEmailService {
	return &syncEmailService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
		maxWorkers:   services.cfg.ConcurrencyConfig.InteractionEventSyncConcurrency,
	}
}

func (s syncEmailService) SyncEmail(ctx context.Context, email model.EmailData) (SyncResult, error) {
	var organizationsData []model.OrganizationData
	var domainsSlice []string
	var name string

	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncEmailService.SyncEmails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "email", email)

	if !s.services.TenantService.Exists(ctx, common.GetTenantFromContext(ctx)) {
		s.log.Errorf("tenant {%s} does not exist", common.GetTenantFromContext(ctx))
		tracing.TraceErr(span, errors.ErrTenantNotValid)
		return SyncResult{}, errors.ErrTenantNotValid
	}

	if email.SentBy != "" {
		domain := utils.ExtractDomain(email.SentBy)
		utils.EnforceSingleValue(domainsSlice, domain)
		name = utils.ExtractName(email.SentBy)
	}

	organizationData := model.OrganizationData{
		BaseData:       email.BaseData,
		Name:           name,
		Domains:        domainsSlice,
		DomainRequired: true,
	}
	organizationsData = append(organizationsData, organizationData)

	orgSyncResult, err := s.services.OrganizationService.SyncOrganizations(ctx, organizationsData)
	if err != nil {
		return orgSyncResult, err
	}

	return orgSyncResult, nil
}
