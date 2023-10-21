package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *LogEntryAggregate) HandleCommand(ctx context.Context, command eventstore.Command) error {
	switch c := command.(type) {
	case *cmd.UpsertLogEntryCommand:
		if c.IsCreateCommand {
			return a.createLogEntry(ctx, c)
		} else {
			return a.updateLogEntry(ctx, c)
		}
	case *cmd.AddTagCommand:
		return a.addTag(ctx, c)
	case *cmd.RemoveTagCommand:
		return a.removeTag(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *LogEntryAggregate) createLogEntry(ctx context.Context, command *cmd.UpsertLogEntryCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.createLogEntry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", command)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(command.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, createdAtNotNil)
	startedAtNotNil := utils.IfNotNilTimeWithDefault(command.DataFields.StartedAt, createdAtNotNil)

	createEvent, err := events.NewLogEntryCreateEvent(a, command.DataFields, command.Source, command.ExternalSystem, createdAtNotNil, updatedAtNotNil, startedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryCreateEvent")
	}
	aggregate.EnrichEventWithMetadata(&createEvent, &span, a.Tenant, command.UserID)

	return a.Apply(createEvent)
}

func (a *LogEntryAggregate) updateLogEntry(ctx context.Context, command *cmd.UpsertLogEntryCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.updateLogEntry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", command)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, utils.Now())
	startedAtNotNil := utils.IfNotNilTimeWithDefault(command.DataFields.StartedAt, a.LogEntry.StartedAt)
	sourceOfTruth := command.Source.SourceOfTruth
	if sourceOfTruth == "" {
		sourceOfTruth = a.LogEntry.Source.SourceOfTruth
	}

	event, err := events.NewLogEntryUpdateEvent(a, command.DataFields.Content, command.DataFields.ContentType,
		sourceOfTruth, updatedAtNotNil, startedAtNotNil, command.DataFields.LoggedOrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryUpdateEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *LogEntryAggregate) addTag(ctx context.Context, command *cmd.AddTagCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.addTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", command)))

	taggedAtNotNil := utils.IfNotNilTimeWithDefault(command.TaggedAt, utils.Now())

	event, err := events.NewLogEntryAddTagEvent(a, command.TagId, taggedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryAddTagEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *LogEntryAggregate) removeTag(ctx context.Context, command *cmd.RemoveTagCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.removeTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", command)))

	event, err := events.NewLogEntryRemoveTagEvent(a, command.TagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryRemoveTagEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}
