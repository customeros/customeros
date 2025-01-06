package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var onboardingStatusByModel = map[model.OnboardingStatus]neo4jenum.OnboardingStatus{
	model.OnboardingStatusNotApplicable: neo4jenum.OnboardingStatusNotApplicable,
	model.OnboardingStatusNotStarted:    neo4jenum.OnboardingStatusNotStarted,
	model.OnboardingStatusOnTrack:       neo4jenum.OnboardingStatusOnTrack,
	model.OnboardingStatusLate:          neo4jenum.OnboardingStatusLate,
	model.OnboardingStatusStuck:         neo4jenum.OnboardingStatusStuck,
	model.OnboardingStatusDone:          neo4jenum.OnboardingStatusDone,
	model.OnboardingStatusSuccessful:    neo4jenum.OnboardingStatusSuccessful,
}

var onboardingStatusByValue = utils.ReverseMap(onboardingStatusByModel)

func MapOnboardingStatusFromModel(input model.OnboardingStatus) neo4jenum.OnboardingStatus {
	return onboardingStatusByModel[input]
}

func MapOnboardingStatusToModel(input neo4jenum.OnboardingStatus) model.OnboardingStatus {
	return onboardingStatusByValue[input]
}
