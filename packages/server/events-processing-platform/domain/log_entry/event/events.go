package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
	"time"
)

const (
	LogEntryCreateV1    = "V1_LOG_ENTRY_CREATE"
	LogEntryUpdateV1    = "V1_LOG_ENTRY_UPDATE"
	LogEntryAddTagV1    = "V1_LOG_ENTRY_ADD_TAG"
	LogEntryRemoveTagV1 = "V1_LOG_ENTRY_REMOVE_TAG"
)

type LogEntryCreateEvent struct {
	Tenant               string                     `json:"tenant" validate:"required"`
	Content              string                     `json:"content"`
	ContentType          string                     `json:"contentType"`
	StartedAt            time.Time                  `json:"startedAt" validate:"required"`
	AuthorUserId         string                     `json:"authorUserId"`
	LoggedOrganizationId string                     `json:"loggedOrganizationId"`
	Source               string                     `json:"source"`
	SourceOfTruth        string                     `json:"sourceOfTruth"`
	AppSource            string                     `json:"appSource"`
	CreatedAt            time.Time                  `json:"createdAt"`
	UpdatedAt            time.Time                  `json:"updatedAt"`
	ExternalSystem       commonmodel.ExternalSystem `json:"externalSystem"`
}

func NewLogEntryCreateEvent(aggregate eventstore.Aggregate, dataFields model.LogEntryDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt, startedAt time.Time) (eventstore.Event, error) {
	eventData := LogEntryCreateEvent{
		Tenant:               aggregate.GetTenant(),
		Content:              dataFields.Content,
		ContentType:          dataFields.ContentType,
		AuthorUserId:         utils.IfNotNilString(dataFields.AuthorUserId),
		LoggedOrganizationId: utils.IfNotNilString(dataFields.LoggedOrganizationId),
		Source:               source.Source,
		SourceOfTruth:        source.SourceOfTruth,
		AppSource:            source.AppSource,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		StartedAt:            startedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate LogEntryCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, LogEntryCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for LogEntryCreateEvent")
	}
	return event, nil
}

type LogEntryUpdateEvent struct {
	Tenant               string    `json:"tenant" validate:"required"`
	Content              string    `json:"content"`
	ContentType          string    `json:"contentType"`
	StartedAt            time.Time `json:"startedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	SourceOfTruth        string    `json:"sourceOfTruth"`
	LoggedOrganizationId string    `json:"loggedOrganizationId"`
}

func NewLogEntryUpdateEvent(aggregate eventstore.Aggregate, content, contentType, sourceOfTruth string, updatedAt, startedAt time.Time, loggedOrganizationId *string) (eventstore.Event, error) {
	eventData := LogEntryUpdateEvent{
		Tenant:               aggregate.GetTenant(),
		Content:              content,
		ContentType:          contentType,
		UpdatedAt:            updatedAt,
		StartedAt:            startedAt,
		SourceOfTruth:        sourceOfTruth,
		LoggedOrganizationId: utils.IfNotNilString(loggedOrganizationId),
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate LogEntryUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, LogEntryUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for LogEntryUpdateEvent")
	}
	return event, nil
}

type LogEntryAddTagEvent struct {
	Tenant   string    `json:"tenant" validate:"required"`
	TagId    string    `json:"tagId" validate:"required"`
	TaggedAt time.Time `json:"taggedAt" validate:"required"`
}

func NewLogEntryAddTagEvent(aggregate eventstore.Aggregate, tagId string, taggedAt time.Time) (eventstore.Event, error) {
	eventData := LogEntryAddTagEvent{
		Tenant:   aggregate.GetTenant(),
		TagId:    tagId,
		TaggedAt: taggedAt,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate LogEntryAddTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, LogEntryAddTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for LogEntryAddTagEvent")
	}
	return event, nil
}

type LogEntryRemoveTagEvent struct {
	Tenant string `json:"tenant" validate:"required"`
	TagId  string `json:"tagId" validate:"required"`
}

func NewLogEntryRemoveTagEvent(aggregate eventstore.Aggregate, tagId string) (eventstore.Event, error) {
	eventData := LogEntryRemoveTagEvent{
		Tenant: aggregate.GetTenant(),
		TagId:  tagId,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate LogEntryRemoveTagEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, LogEntryRemoveTagV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for LogEntryRemoveTagEvent")
	}
	return event, nil
}
