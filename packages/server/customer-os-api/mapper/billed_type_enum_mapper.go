package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var billedTypeByModel = map[model.BilledType]neo4jenum.BilledType{
	model.BilledTypeNone:      neo4jenum.BilledTypeNone,
	model.BilledTypeMonthly:   neo4jenum.BilledTypeMonthly,
	model.BilledTypeQuarterly: neo4jenum.BilledTypeQuarterly,
	model.BilledTypeAnnually:  neo4jenum.BilledTypeAnnually,
	model.BilledTypeOnce:      neo4jenum.BilledTypeOnce,
	model.BilledTypeUsage:     neo4jenum.BilledTypeUsage,
}

var billedTypeByValue = utils.ReverseMap(billedTypeByModel)

func MapBilledTypeFromModel(input model.BilledType) neo4jenum.BilledType {
	return billedTypeByModel[input]
}

func MapBilledTypeToModel(input neo4jenum.BilledType) model.BilledType {
	return billedTypeByValue[input]
}
