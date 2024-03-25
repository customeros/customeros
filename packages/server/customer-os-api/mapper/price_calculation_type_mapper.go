package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var priceCalculationTypeByModel = map[model.CalculationType]neo4jenum.PriceCalculationType{
	model.CalculationTypeRevenueShare: neo4jenum.PriceCalculationTypeRevenueShare,
}

var priceCalculationTypeByValue = utils.ReverseMap(priceCalculationTypeByModel)

func MapPriceCalculationTypeFromModel(input model.CalculationType) neo4jenum.PriceCalculationType {
	return priceCalculationTypeByModel[input]
}

func MapPriceCalculationTypeToModel(input neo4jenum.PriceCalculationType) model.CalculationType {
	return priceCalculationTypeByValue[input]
}
