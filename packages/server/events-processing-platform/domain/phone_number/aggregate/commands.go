package aggregate

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
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
	case *command.FailedPhoneNumberValidationCommand:
		return a.failPhoneNumberValidation(ctx, c)
	case *command.SkippedPhoneNumberValidationCommand:
		return a.skipPhoneNumberValidation(ctx, c)
	case *command.PhoneNumberValidatedCommand:
		return a.phoneNumberValidated(ctx, c)
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

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
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

	event, err := events.NewPhoneNumberUpdateEvent(a, cmd.Tenant, cmd.Source.Source, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberUpdateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) failPhoneNumberValidation(ctx context.Context, cmd *command.FailedPhoneNumberValidationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.failPhoneNumberValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	event, err := events.NewPhoneNumberFailedValidationEvent(a, cmd.Tenant, cmd.RawPhoneNumber, cmd.CountryCodeA2, cmd.ValidationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) skipPhoneNumberValidation(ctx context.Context, cmd *command.SkippedPhoneNumberValidationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.skipPhoneNumberValidation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	event, err := events.NewPhoneNumberSkippedValidationEvent(a, cmd.Tenant, cmd.RawPhoneNumber, cmd.CountryCodeA2, cmd.ValidationSkipReason)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberSkippedValidationEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
	})

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) phoneNumberValidated(ctx context.Context, cmd *command.PhoneNumberValidatedCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.phoneNumberValidated")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.Object("command", cmd))

	event, err := events.NewPhoneNumberValidatedEvent(a, cmd.Tenant, cmd.RawPhoneNumber, cmd.E164, cmd.CountryCodeA2)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberValidatedEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
	})

	return a.Apply(event)
}
