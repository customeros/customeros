package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	entity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/neo4j/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type CountryService interface {
	GetCountryByCodeA3(ctx context.Context, codeA3 string) (*entity.CountryEntity, error)
}

type countryService struct {
	repositories *repository.Repositories
}

func NewCountryService(repositories *repository.Repositories) CountryService {
	return &countryService{
		repositories: repositories,
	}
}

func (s *countryService) GetCountryByCodeA3(ctx context.Context, codeA3 string) (*entity.CountryEntity, error) {
	countryNode, err := s.repositories.CountryRepository.GetCountryByCodeA3(ctx, codeA3)
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToCountryEntity(*countryNode), nil
}

func (s *countryService) mapDbNodeToCountryEntity(node dbtype.Node) *entity.CountryEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.CountryEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CodeA2:    utils.GetStringPropOrEmpty(props, "codeA2"),
		CodeA3:    utils.GetStringPropOrEmpty(props, "codeA3"),
		PhoneCode: utils.GetStringPropOrEmpty(props, "phoneCode"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}
