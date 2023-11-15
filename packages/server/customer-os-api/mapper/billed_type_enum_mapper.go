package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var billedTypeByModel = map[model.BilledType]entity.BilledType{
	model.BilledTypeNone:           entity.BilledTypeNone,
	model.BilledTypeMonthlyBilled:  entity.BilledTypeMonthly,
	model.BilledTypeAnnuallyBilled: entity.BilledTypeAnnually,
	model.BilledTypeOnceBilled:     entity.BilledTypeOnce,
}

var billedTypeByValue = utils.ReverseMap(billedTypeByModel)

func MapBilledTypeFromModel(input model.BilledType) entity.BilledType {
	return billedTypeByModel[input]
}

func MapBilledTypeToModel(input entity.BilledType) model.BilledType {
	return billedTypeByValue[input]
}
