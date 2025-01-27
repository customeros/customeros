package mapper

import (
	postgresEntity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToSku(entity *postgresEntity.SkuEntity) *model.Sku {
	if entity == nil {
		return nil
	}
	output := model.Sku{
		ID:       entity.ID,
		Name:     entity.Name,
		Price:    entity.Price,
		Type:     entity.Type,
		Archived: entity.Archived,
	}
	return &output
}

func MapEntitiesToSkus(entities []*postgresEntity.SkuEntity) []*model.Sku {
	var output []*model.Sku
	for _, v := range entities {
		output = append(output, MapEntityToSku(v))
	}
	return output
}
