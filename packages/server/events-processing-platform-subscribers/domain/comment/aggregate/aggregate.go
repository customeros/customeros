package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	CommentAggregateType eventstore.AggregateType = "comment"
)

type CommentAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Comment *model.Comment
}

func NewCommentAggregateWithTenantAndID(tenant, id string) *CommentAggregate {
	commentAggregate := CommentAggregate{}
	commentAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(CommentAggregateType, tenant, id)
	commentAggregate.SetWhen(commentAggregate.When)
	commentAggregate.Comment = &model.Comment{}
	commentAggregate.Tenant = tenant

	return &commentAggregate
}

func (a *CommentAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.CommentCreateV1:
		return a.onCommentCreate(evt)
	case event.CommentUpdateV1:
		return a.onCommentUpdate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *CommentAggregate) onCommentCreate(evt eventstore.Event) error {
	var eventData event.CommentCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Comment.ID = a.ID
	a.Comment.Tenant = a.Tenant
	a.Comment.Content = eventData.Content
	a.Comment.ContentType = eventData.ContentType
	a.Comment.AuthorUserId = eventData.AuthorUserId
	a.Comment.CommentedIssueId = eventData.CommentedIssueId
	a.Comment.Source = commonmodel.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.Comment.CreatedAt = eventData.CreatedAt
	a.Comment.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Comment.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *CommentAggregate) onCommentUpdate(evt eventstore.Event) error {
	var eventData event.CommentUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Comment.Source.SourceOfTruth = eventData.Source
	}
	if eventData.Source != a.Comment.Source.SourceOfTruth && a.Comment.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Comment.Content == "" {
			a.Comment.Content = eventData.Content
		}
		if a.Comment.ContentType == "" {
			a.Comment.ContentType = eventData.ContentType
		}
	} else {
		a.Comment.Content = eventData.Content
		a.Comment.ContentType = eventData.ContentType
	}
	a.Comment.UpdatedAt = eventData.UpdatedAt
	return nil
}
