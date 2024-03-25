package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var pricingModelByModel = map[model.PricingModel]neo4jenum.PricingModel{
	model.PricingModelSubscription: neo4jenum.PricingModelSubscription,
	model.PricingModelUsage:        neo4jenum.PricingModelUsage,
	model.PricingModelOneTime:      neo4jenum.PricingModelOneTime,
}

var pricingModelByValue = utils.ReverseMap(pricingModelByModel)

func MapPricingModelFromModel(input model.PricingModel) neo4jenum.PricingModel {
	return pricingModelByModel[input]
}

func MapPricingModelToModel(input neo4jenum.PricingModel) model.PricingModel {
	return pricingModelByValue[input]
}
