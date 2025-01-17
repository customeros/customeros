package api_external_system

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"

	cosapi_interfaces "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
)

type externalSystemService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewExternalSystemService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.ExternalSystemService {
	return &externalSystemService{
		log:          log,
		repositories: repositories,
	}
}

func (s *externalSystemService) GetAllExternalSystemInstances(ctx context.Context) (*neo4jentity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetAllExternalSystemInstances")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbRecords, err := s.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(ctx, "Error getting all external system instances", err)
		return nil, err
	}
	var externalSystemEntities neo4jentity.ExternalSystemEntities
	for _, v := range dbRecords {
		externalSystemEntity := neo4jmapper.MapDbNodeToExternalSystem(v)
		externalSystemEntities = append(externalSystemEntities, *externalSystemEntity)
	}
	return &externalSystemEntities, nil
}
