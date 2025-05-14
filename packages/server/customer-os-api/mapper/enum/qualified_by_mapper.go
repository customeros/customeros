package enummapper

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var qualifiedByModel = map[model.QualifiedBy]commonenum.QualifiedBy{
	model.QualifiedBySystem: commonenum.QualifiedBySystem,
	model.QualifiedByUser:   commonenum.QualifiedByUser,
	model.QualifiedByNone:   "",
}

var qualifiedByValue = utils.ReverseMap(qualifiedByModel)

func MapQualifiedByFromModel(input model.QualifiedBy) commonenum.QualifiedBy {
	return qualifiedByModel[input]
}

func MapQualifiedByToModel(input commonenum.QualifiedBy) model.QualifiedBy {
	return qualifiedByValue[input]
}
