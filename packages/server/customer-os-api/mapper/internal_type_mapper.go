package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var internalTypeByModel = map[model.InternalType]neo4jenum.OpportunityInternalType{
	model.InternalTypeNbo:       neo4jenum.OpportunityInternalTypeNBO,
	model.InternalTypeUpsell:    neo4jenum.OpportunityInternalTypeUpsell,
	model.InternalTypeCrossSell: neo4jenum.OpportunityInternalTypeCrossSell,
	model.InternalTypeRenewal:   neo4jenum.OpportunityInternalTypeRenewal,
}

var internalTypeByValue = utils.ReverseMap(internalTypeByModel)

func MapInternalTypeFromModel(input model.InternalType) neo4jenum.OpportunityInternalType {
	return internalTypeByModel[input]
}

func MapInternalTypeToModel(input neo4jenum.OpportunityInternalType) model.InternalType {
	return internalTypeByValue[input]
}
