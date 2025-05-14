package enummapper

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var statusByModel = map[model.QualificationStatus]commonenum.QualificationStatus{
	model.QualificationStatusQualifying:   commonenum.QualificationStatusQualifying,
	model.QualificationStatusQualified:    commonenum.QualificationStatusQualified,
	model.QualificationStatusNotQualified: commonenum.QualificationStatusNotQualified,
	model.QualificationStatusPending:      "",
}

var statusByValue = utils.ReverseMap(statusByModel)

func MapStatusFromModel(input model.QualificationStatus) commonenum.QualificationStatus {
	return statusByModel[input]
}

func MapStatusToModel(input commonenum.QualificationStatus) model.QualificationStatus {
	return statusByValue[input]
}
