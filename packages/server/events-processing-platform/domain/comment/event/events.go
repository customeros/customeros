package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/validator"
	"github.com/pkg/errors"
	"time"
)

const (
	CommentCreateV1 = "V1_COMMENT_CREATE"
	CommentUpdateV1 = "V1_COMMENT_UPDATE"
)

type CommentCreateEvent struct {
	Tenant           string                     `json:"tenant" validate:"required"`
	Content          string                     `json:"content"`
	ContentType      string                     `json:"contentType"`
	AuthorUserId     string                     `json:"authorUserId"`
	CommentedIssueId string                     `json:"commentedIssueId"`
	Source           string                     `json:"source"`
	AppSource        string                     `json:"appSource"`
	CreatedAt        time.Time                  `json:"createdAt"`
	UpdatedAt        time.Time                  `json:"updatedAt"`
	ExternalSystem   commonmodel.ExternalSystem `json:"externalSystem"`
}

func NewCommentCreateEvent(aggregate eventstore.Aggregate, dataFields model.CommentDataFields, source commonmodel.Source, externalSystem commonmodel.ExternalSystem, createdAt, updatedAt time.Time) (eventstore.Event, error) {
	eventData := CommentCreateEvent{
		Tenant:           aggregate.GetTenant(),
		Content:          dataFields.Content,
		ContentType:      dataFields.ContentType,
		AuthorUserId:     utils.IfNotNilString(dataFields.AuthorUserId),
		CommentedIssueId: utils.IfNotNilString(dataFields.CommentedIssueId),
		Source:           source.Source,
		AppSource:        source.AppSource,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate CommentCreateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, CommentCreateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CommentCreateEvent")
	}
	return event, nil
}

type CommentUpdateEvent struct {
	Tenant         string                     `json:"tenant" validate:"required"`
	Content        string                     `json:"content"`
	ContentType    string                     `json:"contentType"`
	UpdatedAt      time.Time                  `json:"updatedAt"`
	Source         string                     `json:"source"`
	ExternalSystem commonmodel.ExternalSystem `json:"externalSystem,omitempty"`
}

func NewCommentUpdateEvent(aggregate eventstore.Aggregate, content, contentType, source string, externalSystem commonmodel.ExternalSystem, updatedAt time.Time) (eventstore.Event, error) {
	eventData := CommentUpdateEvent{
		Tenant:      aggregate.GetTenant(),
		Content:     content,
		ContentType: contentType,
		UpdatedAt:   updatedAt,
		Source:      source,
	}
	if externalSystem.Available() {
		eventData.ExternalSystem = externalSystem
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error validating CommentUpdateEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, CommentUpdateV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for CommentUpdateEvent")
	}
	return event, nil
}
