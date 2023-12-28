package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type UpdateOnboardingStatusEvent struct {
	Tenant             string    `json:"tenant" validate:"required"`
	Status             string    `json:"status" validate:"required"`
	Comments           string    `json:"comments"`
	UpdatedByUserId    string    `json:"updatedByUserId"`
	UpdatedAt          time.Time `json:"updatedAt"`
	CausedByContractId string    `json:"causedByContractId"`
}

func NewUpdateOnboardingStatusEvent(aggregate eventstore.Aggregate, status, comments, updatedByUserId, causedByContractId string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := UpdateOnboardingStatusEvent{
		Tenant:             aggregate.GetTenant(),
		Status:             status,
		Comments:           comments,
		UpdatedByUserId:    updatedByUserId,
		UpdatedAt:          updatedAt,
		CausedByContractId: causedByContractId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate UpdateOnboardingStatusEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationUpdateOnboardingStatusV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for UpdateOnboardingStatusEvent")
	}
	return event, nil
}
