package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
}

type externalSystemService struct {
	log      logger.Logger
	services *Services
}

func NewExternalSystemService(log logger.Logger, services *Services) ExternalSystemService {
	return &externalSystemService{
		log:      log,
		services: services,
	}
}

func (s *externalSystemService) MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.MergeExternalSystem")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("externalSystem", externalSystem))

	if externalSystem == "" {
		return nil
	}

	err := s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(ctx, tenant, externalSystem, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
