package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type PlaceService interface {
	GetForLocation(ctx context.Context, locationId string) (*entity.PlaceEntity, error)
}

type placeService struct {
	repositories *repository.Repositories
}

func NewPlaceService(repositories *repository.Repositories) PlaceService {
	return &placeService{
		repositories: repositories,
	}
}

func (s *placeService) GetForLocation(ctx context.Context, locationId string) (*entity.PlaceEntity, error) {
	dbNodes, err := s.repositories.PlaceRepository.GetAnyForLocation(common.GetTenantFromContext(ctx), locationId)
	if err != nil {
		return nil, err
	}

	if len(dbNodes) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToPlaceEntity(*dbNodes[0]), nil
}

func (s *placeService) mapDbNodeToPlaceEntity(node dbtype.Node) *entity.PlaceEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.PlaceEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
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
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}
