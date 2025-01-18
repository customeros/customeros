package data_fields

import (
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"time"
)

type UserFields struct {
	AppSource                        *string               `json:"appSource,omitempty"`
	Source                           *string               `json:"source,omitempty"`
	ExternalSystem                   *model.ExternalSystem `json:"externalSystem,omitempty"`
	CreatedAt                        *time.Time            `json:"createdAt,omitempty"`
	ShowOnboardingPage               *bool                 `json:"showOnboardingPage,omitempty"`
	OnboardingInboundStepCompleted   *bool                 `json:"onboardingInboundStepCompleted,omitempty"`
	OnboardingOutboundStepCompleted  *bool                 `json:"onboardingOutboundStepCompleted,omitempty"`
	OnboardingCrmStepCompleted       *bool                 `json:"onboardingCrmStepCompleted,omitempty"`
	OnboardingMailstackStepCompleted *bool                 `json:"onboardingMailstackStepCompleted,omitempty"`
	FirstName                        *string               `json:"firstName,omitempty"`
	LastName                         *string               `json:"lastName,omitempty"`
	Name                             *string               `json:"name,omitempty"`
	Timezone                         *string               `json:"timezone,omitempty"`
	ProfilePhotoUrl                  *string               `json:"profilePhotoUrl,omitempty"`
	Internal                         *bool                 `json:"internal,omitempty"`
	Test                             *bool                 `json:"test,omitempty"`
	Bot                              *bool                 `json:"bot,omitempty"`
	Roles                            *[]string             `json:"roles,omitempty"`
}

func (fields UserFields) ExternalSystemAvailable() bool {
	return fields.ExternalSystem != nil && fields.ExternalSystem.Available()
}
