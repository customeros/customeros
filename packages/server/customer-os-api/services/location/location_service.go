package api_location

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-api/entity"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
)

type locationService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories) cosapi_interfaces.LocationService {
	return &locationService{
		log:          log,
		repositories: repositories,
	}
}

func (s *locationService) CreateLocationForEntity(ctx context.Context, entityType commonModel.EntityType, entityId string, source entity.SourceFields) (*neo4jentity.LocationEntity, error) {
	if entityType != commonModel.CONTACT && entityType != commonModel.ORGANIZATION && entityType != commonModel.MEETING {
		return nil, coserrors.ErrInvalidEntityType
	}
	locationNode, err := s.repositories.LocationRepository.CreateLocationForEntity(ctx, common.GetTenantFromContext(ctx), entityType, entityId, source)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToLocationEntity(locationNode), nil
}

func (s *locationService) Update(ctx context.Context, entity neo4jentity.LocationEntity) (*neo4jentity.LocationEntity, error) {
	updatedLocationNode, err := s.repositories.LocationRepository.Update(ctx, common.GetTenantFromContext(ctx), entity)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToLocationEntity(updatedLocationNode), nil
}

func (s *locationService) DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.DetachFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", locationId), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.LocationRepository.RemoveRelationshipAndDeleteOrphans(ctx, entityType, entityId, locationId)

	return err
}
