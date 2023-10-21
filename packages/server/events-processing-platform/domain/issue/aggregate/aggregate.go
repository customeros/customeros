package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	IssueAggregateType eventstore.AggregateType = "issue"
)

type IssueAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Issue *model.Issue
}

func NewIssueAggregateWithTenantAndID(tenant, id string) *IssueAggregate {
	issueAggregate := IssueAggregate{}
	issueAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(IssueAggregateType, tenant, id)
	issueAggregate.SetWhen(issueAggregate.When)
	issueAggregate.Issue = &model.Issue{}
	issueAggregate.Tenant = tenant

	return &issueAggregate
}

func (a *IssueAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.IssueCreateV1:
		return a.onIssueCreate(evt)
	case event.IssueUpdateV1:
		return a.onIssueUpdate(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *IssueAggregate) onIssueCreate(evt eventstore.Event) error {
	var eventData event.IssueCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Issue.ID = a.ID
	a.Issue.Tenant = a.Tenant
	a.Issue.Subject = eventData.Subject
	a.Issue.Description = eventData.Description
	a.Issue.Status = eventData.Status
	a.Issue.Priority = eventData.Priority
	a.Issue.ReportedByOrganization = eventData.ReportedByOrganizationId
	a.Issue.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.Source,
		AppSource:     eventData.AppSource,
	}
	a.Issue.CreatedAt = eventData.CreatedAt
	a.Issue.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Issue.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *IssueAggregate) onIssueUpdate(evt eventstore.Event) error {
	var eventData event.IssueUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if eventData.Source == constants.SourceOpenline {
		a.Issue.Source.SourceOfTruth = eventData.Source
	}
	if eventData.Source != a.Issue.Source.SourceOfTruth && a.Issue.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Issue.Subject == "" {
			a.Issue.Subject = eventData.Subject
		}
		if a.Issue.Description == "" {
			a.Issue.Description = eventData.Description
		}
		if a.Issue.Status == "" {
			a.Issue.Status = eventData.Status
		}
		if a.Issue.Priority == "" {
			a.Issue.Priority = eventData.Priority
		}
	} else {
		a.Issue.Subject = eventData.Subject
		a.Issue.Description = eventData.Description
		a.Issue.Status = eventData.Status
		a.Issue.Priority = eventData.Priority
	}
	a.Issue.UpdatedAt = eventData.UpdatedAt
	return nil
}
