package event

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type ContractUpdateEvent struct {
	Tenant           string                     `json:"tenant" validate:"required"`
	Name             string                     `json:"name"`
	ContractUrl      string                     `json:"contractUrl"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
	ServiceStartedAt *time.Time                 `json:"serviceStartedAt,omitempty"`
	SignedAt         *time.Time                 `json:"signedAt,omitempty"`
	EndedAt          *time.Time                 `json:"endedAt,omitempty"`
	RenewalCycle     string                     `json:"renewalCycle"`
	Status           string                     `json:"status"`
	ExternalSystem   commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
	Source           string                     `json:"source"`
}

func NewContractUpdateEvent(aggr eventstore.Aggregate, dataFields model.ContractDataFields, externalSystem commonmodel.ExternalSystem, source string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := ContractUpdateEvent{
		Tenant:           aggr.GetTenant(),
		Name:             dataFields.Name,
		ContractUrl:      dataFields.ContractUrl,
		ServiceStartedAt: dataFields.ServiceStartedAt,
		SignedAt:         dataFields.SignedAt,
		EndedAt:          dataFields.EndedAt,
		RenewalCycle:     dataFields.RenewalCycle.String(),
		Status:           dataFields.Status.String(),
		UpdatedAt:        updatedAt,
		Source:           source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate ContractUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggr, ContractUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for ContractUpdateEvent")
	}
	return event, nil
}
