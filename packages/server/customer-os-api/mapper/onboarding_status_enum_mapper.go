package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var onboardingStatusByModel = map[model.OnboardingStatus]entity.OnboardingStatus{
	model.OnboardingStatusNotApplicable: entity.OnboardingStatusNotApplicable,
	model.OnboardingStatusNotStarted:    entity.OnboardingStatusNotStarted,
	model.OnboardingStatusOnTrack:       entity.OnboardingStatusOnTrack,
	model.OnboardingStatusLate:          entity.OnboardingStatusLate,
	model.OnboardingStatusStuck:         entity.OnboardingStatusStuck,
	model.OnboardingStatusDone:          entity.OnboardingStatusDone,
	model.OnboardingStatusSuccessful:    entity.OnboardingStatusSuccessful,
}

var onboardingStatusByValue = utils.ReverseMap(onboardingStatusByModel)

func MapOnboardingStatusFromModel(input model.OnboardingStatus) entity.OnboardingStatus {
	return onboardingStatusByModel[input]
}

func MapOnboardingStatusToModel(input entity.OnboardingStatus) model.OnboardingStatus {
	return onboardingStatusByValue[input]
}
