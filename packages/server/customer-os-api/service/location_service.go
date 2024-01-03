package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LocationService interface {
	GetAllForContact(ctx context.Context, contactId string) (*entity.LocationEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*entity.LocationEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*entity.LocationEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.LocationEntities, error)
	CreateLocationForEntity(ctx context.Context, entityType entity.EntityType, entityId string, source entity.SourceFields) (*entity.LocationEntity, error)
	Update(ctx context.Context, entity entity.LocationEntity) (*entity.LocationEntity, error)
	DetachFromEntity(ctx context.Context, entityType entity.EntityType, entityId, locationId string) error
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

func (s *locationService) GetAllForContact(ctx context.Context, contactId string) (*entity.LocationEntities, error) {
	dbNodes, err := s.repositories.LocationRepository.GetAllForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}

	locationEntities := entity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *s.mapDbNodeToLocationEntity(*dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForContacts(ctx context.Context, contactIds []string) (*entity.LocationEntities, error) {
	locations, err := s.repositories.LocationRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	locationEntities := entity.LocationEntities{}
	for _, v := range locations {
		locationEntity := s.mapDbNodeToLocationEntity(*v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganization(ctx context.Context, organizationId string) (*entity.LocationEntities, error) {
	dbNodes, err := s.repositories.LocationRepository.GetAllForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	locationEntities := entity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *s.mapDbNodeToLocationEntity(*dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*entity.LocationEntities, error) {
	locations, err := s.repositories.LocationRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	locationEntities := entity.LocationEntities{}
	for _, v := range locations {
		locationEntity := s.mapDbNodeToLocationEntity(*v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
}

func (s *locationService) CreateLocationForEntity(ctx context.Context, entityType entity.EntityType, entityId string, source entity.SourceFields) (*entity.LocationEntity, error) {
	if entityType != entity.CONTACT && entityType != entity.ORGANIZATION && entityType != entity.MEETING {
		return nil, errors.ErrInvalidEntityType
	}
	locationNode, err := s.repositories.LocationRepository.CreateLocationForEntity(ctx, common.GetTenantFromContext(ctx), entityType, entityId, source)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToLocationEntity(*locationNode), nil
}

func (s *locationService) Update(ctx context.Context, entity entity.LocationEntity) (*entity.LocationEntity, error) {
	updatedLocationNode, err := s.repositories.LocationRepository.Update(ctx, common.GetTenantFromContext(ctx), entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToLocationEntity(*updatedLocationNode), nil
}

func (s *locationService) DetachFromEntity(ctx context.Context, entityType entity.EntityType, entityId, locationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.DetachFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", locationId), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.LocationRepository.RemoveRelationshipAndDeleteOrphans(ctx, entityType, entityId, locationId)

	return err
}

func (s *locationService) mapDbNodeToLocationEntity(node dbtype.Node) *entity.LocationEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.LocationEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Country:       utils.GetStringPropOrEmpty(props, "country"),
		Region:        utils.GetStringPropOrEmpty(props, "region"),
		Locality:      utils.GetStringPropOrEmpty(props, "locality"),
		Address:       utils.GetStringPropOrEmpty(props, "address"),
		Address2:      utils.GetStringPropOrEmpty(props, "address2"),
		Zip:           utils.GetStringPropOrEmpty(props, "zip"),
		AddressType:   utils.GetStringPropOrEmpty(props, "addressType"),
		HouseNumber:   utils.GetStringPropOrEmpty(props, "houseNumber"),
		PostalCode:    utils.GetStringPropOrEmpty(props, "postalCode"),
		PlusFour:      utils.GetStringPropOrEmpty(props, "plusFour"),
		Commercial:    utils.GetBoolPropOrFalse(props, "commercial"),
		Predirection:  utils.GetStringPropOrEmpty(props, "predirection"),
		District:      utils.GetStringPropOrEmpty(props, "district"),
		Street:        utils.GetStringPropOrEmpty(props, "street"),
		RawAddress:    utils.GetStringPropOrEmpty(props, "rawAddress"),
		Latitude:      utils.GetFloatPropOrNil(props, "latitude"),
		Longitude:     utils.GetFloatPropOrNil(props, "longitude"),
		TimeZone:      utils.GetStringPropOrEmpty(props, "timeZone"),
		UtcOffset:     utils.GetInt64PropOrZero(props, "utcOffset"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}
