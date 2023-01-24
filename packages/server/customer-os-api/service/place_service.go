package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type PlaceService interface {
	FindAllForContact(ctx context.Context, contactId string) (*entity.PlaceEntities, error)
	FindAllForOrganization(ctx context.Context, organizationId string) (*entity.PlaceEntities, error)
}

type placeService struct {
	repositories *repository.Repositories
}

func NewPlaceService(repositories *repository.Repositories) PlaceService {
	return &placeService{
		repositories: repositories,
	}
}

func (s *placeService) getDriver() neo4j.Driver {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *placeService) FindAllForContact(ctx context.Context, contactId string) (*entity.PlaceEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.AddressRepository.FindAllForContact(session, common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}

	addressEntities := entity.PlaceEntities{}
	for _, dbNode := range dbNodes {
		addressEntities = append(addressEntities, *s.mapDbNodeToAddressEntity(dbNode))
	}
	return &addressEntities, nil
}

func (s *placeService) FindAllForOrganization(ctx context.Context, organizationId string) (*entity.PlaceEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbNodes, err := s.repositories.AddressRepository.FindAllForOrganization(session, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	addressEntities := entity.PlaceEntities{}
	for _, dbNode := range dbNodes {
		addressEntities = append(addressEntities, *s.mapDbNodeToAddressEntity(dbNode))
	}
	return &addressEntities, nil
}

func (s *placeService) mapDbNodeToAddressEntity(node *dbtype.Node) *entity.PlaceEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.PlaceEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		Country:       utils.GetStringPropOrEmpty(props, "country"),
		State:         utils.GetStringPropOrEmpty(props, "state"),
		City:          utils.GetStringPropOrEmpty(props, "city"),
		Address:       utils.GetStringPropOrEmpty(props, "address"),
		Address2:      utils.GetStringPropOrEmpty(props, "address2"),
		Zip:           utils.GetStringPropOrEmpty(props, "zip"),
		Phone:         utils.GetStringPropOrEmpty(props, "phone"),
		Fax:           utils.GetStringPropOrEmpty(props, "fax"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &result
}
