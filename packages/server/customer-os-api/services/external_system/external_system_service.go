package api_external_system

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
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
	span, ctx := telemetry.StartServiceSpan(ctx, "ExternalSystemService.GetAllExternalSystemInstances")
	defer span.Finish()

	dbRecords, err := s.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		span.TraceError(err)
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
