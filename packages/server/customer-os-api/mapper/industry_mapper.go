package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToIndustry(entity *neo4jentity.IndustryEntity) *model.Industry {
	industry := model.Industry{
		Code: entity.Code,
		Name: entity.Name,
	}
	return &industry
}

func MapEntitiesToIndustries(entities *neo4jentity.IndustryEntities) []*model.Industry {
	var industries []*model.Industry
	for _, industryEntity := range *entities {
		industries = append(industries, MapEntityToIndustry(&industryEntity))
	}
	return industries
}
