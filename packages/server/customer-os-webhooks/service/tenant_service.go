package service

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/caches"
	"github.com/customeros/customeros/packages/server/customer-os-webhooks/repository"
)

type TenantService interface {
	Exists(ctx context.Context, tenant string) bool
}

type tenantService struct {
	log    logger.Logger
	neo4j  *repository.Repositories
	caches *caches.Cache
}

func NewTenantService(log logger.Logger, repositories *repository.Repositories, caches *caches.Cache) TenantService {
	return &tenantService{
		log:    log,
		neo4j:  repositories,
		caches: caches,
	}
}

func (s *tenantService) Exists(ctx context.Context, tenant string) bool {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantService.Exists")
	defer spans.Finish()

	if !s.caches.CheckTenant(tenant) {
		_, err := s.neo4j.TenantRepository.GetTenant(ctx, tenant)
		if err != nil {
			spans.LogKV("output", false)
			return false
		}
		s.caches.AddTenant(tenant)
	}
	spans.LogKV("output", true)
	return true
}
