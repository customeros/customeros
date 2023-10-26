package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"time"
)

const (
	InteractionEventCreateV1             = "V1_INTERACTION_EVENT_CREATE"
	InteractionEventUpdateV1             = "V1_INTERACTION_EVENT_UPDATE"
	InteractionEventRequestSummaryV1     = "V1_INTERACTION_EVENT_REQUEST_SUMMARY"
	InteractionEventReplaceSummaryV1     = "V1_INTERACTION_EVENT_REPLACE_SUMMARY"
	InteractionEventRequestActionItemsV1 = "V1_INTERACTION_EVENT_REQUEST_ACTION_ITEMS"
	InteractionEventReplaceActionItemsV1 = "V1_INTERACTION_EVENT_REPLACE_ACTION_ITEMS"
)

type InteractionEventCreateEvent struct {
	Tenant          string                `json:"tenant" validate:"required"`
	Content         string                `json:"content"`
	ContentType     string                `json:"contentType"`
	Channel         string                `json:"channel"`
	ChannelData     string                `json:"channelData"`
	EventType       string                `json:"eventType"`
	Identifier      string                `json:"identifier"`
	PartOfIssueId   string                `json:"partOfIssueId,omitempty" validate:"required_without=PartOfSessionId"`
	PartOfSessionId string                `json:"partOfSessionId,omitempty" validate:"required_without=PartOfIssueId"`
	Source          string                `json:"source"`
	AppSource       string                `json:"appSource"`
	CreatedAt       time.Time             `json:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	ExternalSystem  cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	Hide            bool                  `json:"hide"`
}

func NewInteractionEventCreateEvent(aggregate eventstore.Aggregate, dataFields model.InteractionEventDataFields, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InteractionEventCreateEvent{
		Tenant:          aggregate.GetTenant(),
		Content:         dataFields.Content,
		ContentType:     dataFields.ContentType,
		Channel:         dataFields.Channel,
		ChannelData:     dataFields.ChannelData,
		EventType:       dataFields.EventType,
		Identifier:      dataFields.Identifier,
		Hide:            dataFields.Hide,
		PartOfIssueId:   utils.IfNotNilString(dataFields.PartOfIssueId),
		PartOfSessionId: utils.IfNotNilString(dataFields.PartOfSessionId),
		Source:          source.Source,
		AppSource:       source.AppSource,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type InteractionEventUpdateEvent struct {
	Tenant         string                `json:"tenant" validate:"required"`
	Content        string                `json:"content"`
	ContentType    string                `json:"contentType"`
	Channel        string                `json:"channel"`
	ChannelData    string                `json:"channelData"`
	EventType      string                `json:"eventType"`
	Identifier     string                `json:"identifier"`
	UpdatedAt      time.Time             `json:"updatedAt"`
	Source         string                `json:"source"`
	ExternalSystem cmnmod.ExternalSystem `json:"externalSystem,omitempty"`
	Hide           bool                  `json:"hide"`
}

func NewInteractionEventUpdateEvent(aggregate eventstore.Aggregate, dataFields model.InteractionEventDataFields, source string, externalSystem cmnmod.ExternalSystem, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InteractionEventUpdateEvent{
		Tenant:      aggregate.GetTenant(),
		Content:     dataFields.Content,
		ContentType: dataFields.ContentType,
		Channel:     dataFields.Channel,
		ChannelData: dataFields.ChannelData,
		EventType:   dataFields.EventType,
		Identifier:  dataFields.Identifier,
		Hide:        dataFields.Hide,
		UpdatedAt:   updatedAt,
		Source:      source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type InteractionEventRequestSummaryEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewInteractionEventRequestSummaryEvent(aggregate eventstore.Aggregate, tenant string) (eventstore.Event, error) {
	eventData := InteractionEventRequestSummaryEvent{
		Tenant:      tenant,
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventRequestSummaryV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type InteractionEventReplaceSummaryEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	Summary     string    `json:"summary"`
	ContentType string    `json:"contentType"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewInteractionEventReplaceSummaryEvent(aggregate eventstore.Aggregate, tenant, summary, contentType string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InteractionEventReplaceSummaryEvent{
		Tenant:      tenant,
		Summary:     summary,
		UpdatedAt:   updatedAt,
		ContentType: contentType,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventReplaceSummaryV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type InteractionEventRequestActionsItemsEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	RequestedAt time.Time `json:"requestedAt"`
}

func NewInteractionEventRequestActionItemsEvent(aggregate eventstore.Aggregate, tenant string) (eventstore.Event, error) {
	eventData := InteractionEventRequestActionsItemsEvent{
		Tenant:      tenant,
		RequestedAt: utils.Now(),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventRequestActionItemsV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}

type InteractionEventReplaceActionItemsEvent struct {
	Tenant      string    `json:"tenant" validate:"required"`
	ActionItems []string  `json:"actionItems"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewInteractionEventReplaceActionItemsEvent(aggregate eventstore.Aggregate, tenant string, actionItems []string, updatedAt time.Time) (eventstore.Event, error) {
	eventData := InteractionEventReplaceActionItemsEvent{
		Tenant:      tenant,
		ActionItems: actionItems,
		UpdatedAt:   updatedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, err
	}

	event := eventstore.NewBaseEvent(aggregate, InteractionEventReplaceActionItemsV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, err
	}
	return event, nil
}
