package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToOffering(entity *neo4jentity.OfferingEntity) *model.Offering {
	if entity == nil {
		return nil
	}
	return &model.Offering{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		Name:                  entity.Name,
		Active:                entity.Active,
		Type:                  utils.ToPtr(MapOfferingTypeToModel(entity.Type)),
		PricingModel:          utils.ToPtr(MapPricingModelToModel(entity.PricingModel)),
		PricingPeriodInMonths: entity.PricingPeriodInMonths,
		Currency:              utils.ToPtr(MapCurrencyToModel(entity.Currency)),
		DefaultPrice:          entity.Price,
		PriceCalculated:       entity.PriceCalculated,
		Conditional:           entity.Conditional,
		Taxable:               entity.Taxable,
		PriceCalculation: &model.PriceCalculation{
			CalculationType:        utils.ToPtr(MapPriceCalculationTypeToModel(entity.PriceCalculation.Type)),
			RevenueSharePercentage: entity.PriceCalculation.RevenueSharePercentage,
		},
		Conditionals: &model.Conditionals{
			MinimumChargePeriod: utils.ToPtr(MapChargePeriodToModel(entity.Conditionals.MinimumChargePeriod)),
			MinimumChargeAmount: entity.Conditionals.MinimumChargeAmount,
		},
	}
}

func MapEntitiesToOfferings(entities *neo4jentity.OfferingEntities) []*model.Offering {
	var models []*model.Offering
	for _, entity := range *entities {
		models = append(models, MapEntityToOffering(&entity))
	}
	return models
}
