package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
)

type TenantService interface {
	Exists(ctx context.Context, tenant string) bool
}

type tenantService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewTenantService(log logger.Logger, repositories *repository.Repositories) TenantService {
	return &tenantService{
		log:          log,
		repositories: repositories,
	}
}

func (s *tenantService) Exists(ctx context.Context, tenant string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.Exists")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	_, err := s.repositories.TenantRepository.GetTenant(ctx, tenant)
	if err != nil {
		return false
	}
	return true
}
