package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	MergeExternalSystem(ctx context.Context, tenant, externalSystem string) error
	SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith LinkWith) error
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

func (s *externalSystemService) SetPrimaryExternalId(ctx context.Context, externalSystem, externalId string, linkWith LinkWith) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.SetPrimaryExternalId")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	// Create external system if it doesn't exist
	err := s.MergeExternalSystem(ctx, tenant, externalSystem)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Set primary external id
	err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.SetPrimaryExternalId(ctx, tenant, externalSystem, externalId, linkWith.Type.Neo4jLabel(), linkWith.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Send completion event if link with is an organization
	if linkWith.Type == model.ORGANIZATION {
		utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}
	return nil
}
