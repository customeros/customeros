package event

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

type InteractionSessionCreateEvent struct {
	Tenant         string                `json:"tenant" validate:"required"`
	Channel        string                `json:"channel"`
	ChannelData    string                `json:"channelData"`
	Type           string                `json:"type"`
	Identifier     string                `json:"identifier"`
	Name           string                `json:"name"`
	Status         string                `json:"status"`
	Source         string                `json:"source"`
	AppSource      string                `json:"appSource"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	ExternalSystem cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewInteractionSessionCreateEvent(aggregate eventstore.Aggregate, dataFields model.InteractionSessionDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InteractionSessionCreateEvent{
		Tenant:      aggregate.GetTenant(),
		Channel:     dataFields.Channel,
		ChannelData: dataFields.ChannelData,
		Identifier:  dataFields.Identifier,
		Name:        dataFields.Name,
		Status:      dataFields.Status,
		Type:        dataFields.Type,
		Source:      source.Source,
		AppSource:   source.AppSource,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate InteractionSessionCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionSessionCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for InteractionSessionCreateEvent")
	}
	return event, nil
}
