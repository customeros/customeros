package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	LogEntryAggregateType eventstore.AggregateType = "log_entry"
)

type LogEntryAggregate struct {
	*aggregate.CommonTenantIdAggregate
	LogEntry *model.LogEntry
}

func NewLogEntryAggregateWithTenantAndID(tenant, id string) *LogEntryAggregate {
	logEntryAggregate := LogEntryAggregate{}
	logEntryAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(LogEntryAggregateType, tenant, id)
	logEntryAggregate.SetWhen(logEntryAggregate.When)
	logEntryAggregate.LogEntry = &model.LogEntry{}
	logEntryAggregate.Tenant = tenant

	return &logEntryAggregate
}

func (a *LogEntryAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.LogEntryCreateV1:
		return a.onLogEntryCreate(evt)
	case event.LogEntryUpdateV1:
		return a.onLogEntryUpdate(evt)
	case event.LogEntryAddTagV1:
		return a.onLogEntryAddTag(evt)
	case event.LogEntryRemoveTagV1:
		return a.onLogEntryRemoveTag(evt)
	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = evt.GetEventType()
		return err
	}
}

func (a *LogEntryAggregate) onLogEntryCreate(evt eventstore.Event) error {
	var eventData event.LogEntryCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.LogEntry.ID = a.ID
	a.LogEntry.Tenant = a.Tenant
	a.LogEntry.Content = eventData.Content
	a.LogEntry.ContentType = eventData.ContentType
	a.LogEntry.AuthorUserId = eventData.AuthorUserId
	if eventData.LoggedOrganizationId != "" {
		a.LogEntry.LoggedOrganizationIds = utils.AddToListIfNotExists(a.LogEntry.LoggedOrganizationIds, eventData.LoggedOrganizationId)
	}
	a.LogEntry.StartedAt = eventData.StartedAt
	a.LogEntry.Source = commonmodel.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.LogEntry.CreatedAt = eventData.CreatedAt
	a.LogEntry.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.LogEntry.ExternalSystems = []commonmodel.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *LogEntryAggregate) onLogEntryUpdate(evt eventstore.Event) error {
	var eventData event.LogEntryUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.LogEntry.Content = eventData.Content
	a.LogEntry.ContentType = eventData.ContentType
	a.LogEntry.StartedAt = eventData.StartedAt
	a.LogEntry.UpdatedAt = eventData.UpdatedAt
	if eventData.LoggedOrganizationId != "" {
		a.LogEntry.LoggedOrganizationIds = utils.AddToListIfNotExists(a.LogEntry.LoggedOrganizationIds, eventData.LoggedOrganizationId)
	}
	if eventData.SourceOfTruth != "" {
		a.LogEntry.Source.SourceOfTruth = eventData.SourceOfTruth
	}
	return nil
}

func (a *LogEntryAggregate) onLogEntryAddTag(evt eventstore.Event) error {
	var eventData event.LogEntryAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.LogEntry.TagIds = append(a.LogEntry.TagIds, eventData.TagId)
	a.LogEntry.TagIds = utils.RemoveDuplicates(a.LogEntry.TagIds)

	return nil
}

func (a *LogEntryAggregate) onLogEntryRemoveTag(evt eventstore.Event) error {
	var eventData event.LogEntryRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.LogEntry.TagIds = utils.RemoveFromList(a.LogEntry.TagIds, eventData.TagId)

	return nil
}
