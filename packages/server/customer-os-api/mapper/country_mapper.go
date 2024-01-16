package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToCountry(entity *entity.CountryEntity) *model.Country {
	if entity == nil {
		return nil
	}
	return &model.Country{
		Name:      entity.Name,
		CodeA2:    entity.CodeA2,
		CodeA3:    entity.CodeA3,
		PhoneCode: entity.PhoneCode,
	}
}

func MapEntitiesToCountries(entities *entity.CountryEntities) []*model.Country {
	var countries []*model.Country
	for _, entity := range *entities {
		countries = append(countries, MapEntityToCountry(&entity))
	}
	return countries
}
