package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func (a *InteractionSessionAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpsertInteractionSessionCommand:
		if c.IsCreateCommand {
			return a.createInteractionSession(ctx, c)
		}
		return nil
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *InteractionSessionAggregate) createInteractionSession(ctx context.Context, cmd *command.UpsertInteractionSessionCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InteractionSessionAggregate.createInteractionSession")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	createEvent, err := event.NewInteractionSessionCreateEvent(a, cmd.DataFields, cmd.Source, cmd.ExternalSystem, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.GetLoggedInUserId(),
		App:    cmd.GetAppSource(),
	})

	return a.Apply(createEvent)
}
