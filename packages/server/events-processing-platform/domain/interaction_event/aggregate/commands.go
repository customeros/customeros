package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *InteractionEventAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	switch c := cmd.(type) {
	case *command.UpsertInteractionEventCommand:
		if c.IsCreateCommand {
			return a.createInteractionEvent(ctx, c)
		} else {
			return a.updateInteractionEvent(ctx, c)
		}
	default:
		return eventstore.ErrInvalidCommandType
	}
}

func (a *InteractionEventAggregate) createInteractionEvent(ctx context.Context, cmd *command.UpsertInteractionEventCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.createInteractionEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := event.NewInteractionEventCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(createEvent)
}

func (a *InteractionEventAggregate) updateInteractionEvent(ctx context.Context, cmd *command.UpsertInteractionEventCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.updateInteractionEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	source := cmd.Source.Source
	if source == "" {
		source = a.InteractionEvent.Source.Source
	}

	updateEvent, err := event.NewInteractionEventUpdateEvent(a, cmd.DataFields, cmd.Source.Source, cmd.ExternalSystem, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *InteractionEventAggregate) RequestSummary(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.RequestSummary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := event.NewInteractionEventRequestSummaryEvent(a, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventRequestSummaryEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) ReplaceSummary(ctx context.Context, tenant, summary, contentType string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.ReplaceSummary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())

	event, err := event.NewInteractionEventReplaceSummaryEvent(a, tenant, summary, contentType, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventReplaceSummaryEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) RequestActionItems(ctx context.Context, tenant string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.RequestActionItems")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := event.NewInteractionEventRequestActionItemsEvent(a, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventRequestActionItemsEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *InteractionEventAggregate) ReplaceActionItems(ctx context.Context, tenant string, actionItems []string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionEventAggregate.ReplaceActionItems")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())

	event, err := event.NewInteractionEventReplaceActionItemsEvent(a, tenant, actionItems, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewInteractionEventReplaceActionItemsEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
