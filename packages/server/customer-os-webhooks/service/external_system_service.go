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

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
}

type externalSystemService struct {
	log          logger.Logger
	repositories *repository.Repositories
	caches       *caches.Cache
}

func NewExternalSystemService(log logger.Logger, repositories *repository.Repositories, caches *caches.Cache) ExternalSystemService {
	return &externalSystemService{
		log:          log,
		repositories: repositories,
		caches:       caches,
	}
}

func (s *externalSystemService) MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.MergeExternalSystem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem))

	if !s.caches.CheckExternalSystem(tenant, externalSystem) {
		err := s.repositories.ExternalSystemRepository.MergeExternalSystem(ctx, tenant, externalSystem, externalSystem)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		s.caches.AddExternalSystem(tenant, externalSystem)
	}
	return nil
}
