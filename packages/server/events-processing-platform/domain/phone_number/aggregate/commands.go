package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *PhoneNumberAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.UpsertPhoneNumberCommand:
		if c.IsCreateCommand {
			return a.createPhoneNumber(ctx, c)
		} else {
			return a.updatePhoneNumber(ctx, c)
		}
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *PhoneNumberAggregate) createPhoneNumber(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.createPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	cmd.Source.SetDefaultValues()

	event, err := events.NewPhoneNumberCreateEvent(a, cmd.Tenant, cmd.RawPhoneNumber, cmd.Source, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberCreateEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) updatePhoneNumber(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.updatePhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("command", fmt.Sprintf("%+v", cmd)))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	if cmd.Source.Source == "" {
		cmd.Source.Source = constants.SourceOpenline
	}

	event, err := events.NewPhoneNumberUpdateEvent(a, cmd.Tenant, cmd.Source.Source, cmd.RawPhoneNumber, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberUpdateEvent")
	}

	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(event)
}
