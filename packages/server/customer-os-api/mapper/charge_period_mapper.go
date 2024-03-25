package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var chargePeriodByModel = map[model.ChargePeriod]neo4jenum.ChargePeriod{
	model.ChargePeriodMonthly:   neo4jenum.ChargePeriodMonthly,
	model.ChargePeriodQuarterly: neo4jenum.ChargePeriodQuarterly,
	model.ChargePeriodAnnually:  neo4jenum.ChargePeriodAnnually,
}

var chargePeriodByValue = utils.ReverseMap(chargePeriodByModel)

func MapChargePeriodFromModel(input model.ChargePeriod) neo4jenum.ChargePeriod {
	return chargePeriodByModel[input]
}

func MapChargePeriodToModel(input neo4jenum.ChargePeriod) model.ChargePeriod {
	return chargePeriodByValue[input]
}
