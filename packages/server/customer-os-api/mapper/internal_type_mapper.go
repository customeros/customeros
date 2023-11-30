package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var internalTypeByModel = map[model.InternalType]entity.InternalType{
	model.InternalTypeNbo:       entity.InternalTypeNbo,
	model.InternalTypeUpsell:    entity.InternalTypeUpsell,
	model.InternalTypeCrossSell: entity.InternalTypeCrossSell,
	model.InternalTypeRenewal:   entity.InternalTypeRenewal,
}

var internalTypeByValue = utils.ReverseMap(internalTypeByModel)

func MapInternalTypeFromModel(input model.InternalType) entity.InternalType {
	return internalTypeByModel[input]
}

func MapInternalTypeToModel(input entity.InternalType) model.InternalType {
	return internalTypeByValue[input]
}
