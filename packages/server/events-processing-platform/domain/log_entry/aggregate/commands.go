package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *LogEntryAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpsertLogEntryCommand:
		if c.IsCreateCommand {
			return a.createLogEntry(ctx, c)
		} else {
			return a.updateLogEntry(ctx, c)
		}
	case *command.AddTagCommand:
		return a.addTag(ctx, c)
	case *command.RemoveTagCommand:
		return a.removeTag(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *LogEntryAggregate) createLogEntry(ctx context.Context, cmd *command.UpsertLogEntryCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.createLogEntry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	startedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.DataFields.StartedAt, createdAtNotNil)

	createEvent, err := event.NewLogEntryCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil, startedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *LogEntryAggregate) updateLogEntry(ctx context.Context, cmd *command.UpsertLogEntryCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.updateLogEntry")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	startedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.DataFields.StartedAt, a.LogEntry.StartedAt)
	sourceOfTruth := cmd.Source.SourceOfTruth
	if sourceOfTruth == "" {
		sourceOfTruth = a.LogEntry.Source.SourceOfTruth
	}

	updateEvent, err := event.NewLogEntryUpdateEvent(a, cmd.DataFields.Content, cmd.DataFields.ContentType,
		sourceOfTruth, updatedAtNotNil, startedAtNotNil, cmd.DataFields.LoggedOrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *LogEntryAggregate) addTag(ctx context.Context, cmd *command.AddTagCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.addTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	taggedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.TaggedAt, utils.Now())

	addTagEvent, err := event.NewLogEntryAddTagEvent(a, cmd.TagId, taggedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryAddTagEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addTagEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    "", // TODO add appSource into grpc message
	})

	return a.Apply(addTagEvent)
}

func (a *LogEntryAggregate) removeTag(ctx context.Context, cmd *command.RemoveTagCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LogEntryAggregate.removeTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	removeTagEvent, err := event.NewLogEntryRemoveTagEvent(a, cmd.TagId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLogEntryRemoveTagEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&removeTagEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    "", // TODO add appSource into grpc message
	})

	return a.Apply(removeTagEvent)
}
