package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var billedTypeByModel = map[model.BilledType]entity.BilledType{
	model.BilledTypeNone:      entity.BilledTypeNone,
	model.BilledTypeMonthly:   entity.BilledTypeMonthly,
	model.BilledTypeQuarterly: entity.BilledTypeQuarterly,
	model.BilledTypeAnnually:  entity.BilledTypeAnnually,
	model.BilledTypeOnce:      entity.BilledTypeOnce,
	model.BilledTypeUsage:     entity.BilledTypeUsage,
}

var billedTypeByValue = utils.ReverseMap(billedTypeByModel)

func MapBilledTypeFromModel(input model.BilledType) entity.BilledType {
	return billedTypeByModel[input]
}

func MapBilledTypeToModel(input entity.BilledType) model.BilledType {
	return billedTypeByValue[input]
}
