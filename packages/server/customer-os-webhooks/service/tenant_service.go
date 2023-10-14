package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantService interface {
	Exists(ctx context.Context, tenant string) bool
}

type tenantService struct {
	log          logger.Logger
	repositories *repository.Repositories
	caches       *caches.Cache
}

func NewTenantService(log logger.Logger, repositories *repository.Repositories, caches *caches.Cache) TenantService {
	return &tenantService{
		log:          log,
		repositories: repositories,
		caches:       caches,
	}
}

func (s *tenantService) Exists(ctx context.Context, tenant string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantService.Exists")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if !s.caches.CheckTenant(tenant) {
		_, err := s.repositories.TenantRepository.GetTenant(ctx, tenant)
		if err != nil {
			span.LogFields(log.Bool("output", false))
			return false
		}
		s.caches.AddTenant(tenant)
	}
	span.LogFields(log.Bool("output", true))
	return true
}
