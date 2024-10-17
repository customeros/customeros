package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
)

type CustomFieldTemplateService interface {
	GetAll(ctx context.Context) (*neo4jentity.CustomFieldTemplateEntities, error)
	//Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId, opportunityId *string, input *repository.OpportunitySaveFields) (*string, error)
}

type customFieldTemplateService struct {
	log      logger.Logger
	services *Services
}

func NewCustomFieldTemplateService(log logger.Logger, services *Services) CustomFieldTemplateService {
	return &customFieldTemplateService{
		log:      log,
		services: services,
	}
}

func (s *customFieldTemplateService) GetAll(ctx context.Context) (*neo4jentity.CustomFieldTemplateEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateService.GetAll")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	dbNodes, err := s.services.Neo4jRepositories.CustomFieldTemplateReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	customFieldTemplateEntities := neo4jentity.CustomFieldTemplateEntities{}
	for _, dbNode := range dbNodes {
		customFieldTemplateEntities = append(customFieldTemplateEntities, *neo4jmapper.MapDbNodeToCustomFieldTemplateEntity(dbNode))
	}
	return &customFieldTemplateEntities, nil
}
