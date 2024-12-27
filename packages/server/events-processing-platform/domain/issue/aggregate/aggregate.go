package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"strings"
)

const (
	IssueAggregateType eventstore.AggregateType = "issue"
)

type IssueAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Issue *model.Issue
}

func NewIssueAggregateWithTenantAndID(tenant, id string) *IssueAggregate {
	issueAggregate := IssueAggregate{}
	issueAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(IssueAggregateType, tenant, id)
	issueAggregate.SetWhen(issueAggregate.When)
	issueAggregate.Issue = &model.Issue{}
	issueAggregate.Tenant = tenant

	return &issueAggregate
}

func (a *IssueAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.IssueAddUserAssigneeV1:
		return a.onIssueAddUserAssignee(evt)
	case event.IssueRemoveUserAssigneeV1:
		return a.onIssueRemoveUserAssignee(evt)
	case event.IssueAddUserFollowerV1:
		return a.onIssueAddUserFollower(evt)
	case event.IssueRemoveUserFollowerV1:
		return a.onIssueRemoveUserFollower(evt)
	default:
		if strings.HasPrefix(evt.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *IssueAggregate) onIssueAddUserAssignee(evt eventstore.Event) error {
	var eventData event.IssueAddUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.AddAssignedToUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueRemoveUserAssignee(evt eventstore.Event) error {
	var eventData event.IssueRemoveUserAssigneeEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.RemoveAssignedToUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueAddUserFollower(evt eventstore.Event) error {
	var eventData event.IssueAddUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.AddFollowedByUserId(eventData.UserId)
	return nil
}

func (a *IssueAggregate) onIssueRemoveUserFollower(evt eventstore.Event) error {
	var eventData event.IssueRemoveUserFollowerEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.RemoveFollowedByUserId(eventData.UserId)
	return nil
}
