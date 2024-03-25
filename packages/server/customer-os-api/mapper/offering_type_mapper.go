package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var offeringTypeByModel = map[model.OfferingType]neo4jenum.OfferingType{
	model.OfferingTypeProduct: neo4jenum.OfferingTypeProduct,
	model.OfferingTypeService: neo4jenum.OfferingTypeService,
}

var offeringTypeByValue = utils.ReverseMap(offeringTypeByModel)

func MapOfferingTypeFromModel(input model.OfferingType) neo4jenum.OfferingType {
	return offeringTypeByModel[input]
}

func MapOfferingTypeToModel(input neo4jenum.OfferingType) model.OfferingType {
	return offeringTypeByValue[input]
}
