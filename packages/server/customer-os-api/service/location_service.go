package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationService interface {
	GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error)
	CreateLocationForEntity(ctx context.Context, entityType commonModel.EntityType, entityId string, source entity.SourceFields) (*neo4jentity.LocationEntity, error)
	Update(ctx context.Context, entity neo4jentity.LocationEntity) (*neo4jentity.LocationEntity, error)
	DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, locationId string) error
}

type locationService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewLocationService(log logger.Logger, repositories *repository.Repositories) LocationService {
	return &locationService{
		log:          log,
		repositories: repositories,
	}
}

func (s *locationService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *locationService) GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error) {
	dbNodes, err := s.repositories.LocationRepository.GetAllForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}

	locationEntities := neo4jentity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *neo4jmapper.MapDbNodeToLocationEntity(dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error) {
	locations, err := s.repositories.LocationRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	locationEntities := neo4jentity.LocationEntities{}
	for _, v := range locations {
		locationEntity := neo4jmapper.MapDbNodeToLocationEntity(v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error) {
	dbNodes, err := s.repositories.LocationRepository.GetAllForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	locationEntities := neo4jentity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *neo4jmapper.MapDbNodeToLocationEntity(dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error) {
	locations, err := s.repositories.LocationRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	locationEntities := neo4jentity.LocationEntities{}
	for _, v := range locations {
		locationEntity := neo4jmapper.MapDbNodeToLocationEntity(v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
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
