package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

var icpFitByModel = map[model.IcpFit]enum.IcpFit{
	model.IcpFitIcpFit:    enum.IcpIsFit,
	model.IcpFitIcpNotFit: enum.IcpNotFit,
	model.IcpFitIcpNotSet: enum.IcpNotSet,
}

var icpFitByValue = utils.ReverseMap(icpFitByModel)

func MapIcpFitFromModel(input model.IcpFit) enum.IcpFit {
	return icpFitByModel[input]
}

func MapIcpFitToModel(input enum.IcpFit) model.IcpFit {
	return icpFitByValue[input]
}
