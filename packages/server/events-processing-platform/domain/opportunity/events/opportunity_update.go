package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	opportunityevent "github.com/openline-ai/openline-customer-os/packages/server/events/event/opportunity"
	"time"
)

type OpportunityUpdateEvent struct {
	Tenant            string                     `json:"tenant" validate:"required"`
	Name              string                     `json:"name"`
	Amount            float64                    `json:"amount"`
	MaxAmount         float64                    `json:"maxAmount"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
	OwnerUserId       string                     `json:"ownerUserId"`
	InternalStage     string                     `json:"internalStage"`
	Source            string                     `json:"source"`
	AppSource         string                     `json:"appSource"`
	ExternalSystem    commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	ExternalStage     string                     `json:"externalStage"`
	ExternalType      string                     `json:"externalType"`
	EstimatedClosedAt *time.Time                 `json:"estimatedClosedAt,omitempty"`
	FieldsMask        []string                   `json:"fieldsMask"`
	Currency          string                     `json:"currency"`
	NextSteps         string                     `json:"nextSteps"`
	LikelihoodRate    int64                      `json:"likelihoodRate"`
}

func (e OpportunityUpdateEvent) UpdateName() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskName)
}

func (e OpportunityUpdateEvent) UpdateAmount() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskAmount)
}

func (e OpportunityUpdateEvent) UpdateMaxAmount() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskMaxAmount)
}

func (e OpportunityUpdateEvent) UpdateExternalStage() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskExternalStage)
}

func (e OpportunityUpdateEvent) UpdateExternalType() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskExternalType)
}

func (e OpportunityUpdateEvent) UpdateEstimatedClosedAt() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskEstimatedClosedAt)
}

func (e OpportunityUpdateEvent) UpdateOwnerUserId() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskOwnerUserId)
}

func (e OpportunityUpdateEvent) UpdateInternalStage() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskInternalStage) && e.InternalStage != ""
}

func (e OpportunityUpdateEvent) UpdateCurrency() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskCurrency)
}

func (e OpportunityUpdateEvent) UpdateNextSteps() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskNextSteps)
}

func (e OpportunityUpdateEvent) UpdateLikelihoodRate() bool {
	return utils.Contains(e.FieldsMask, opportunityevent.FieldMaskLikelihoodRate)
}
