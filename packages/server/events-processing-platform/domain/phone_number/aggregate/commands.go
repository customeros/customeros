package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
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
	switch c := cmd.(type) {
	case *command.UpsertPhoneNumberCommand:
		if c.IsCreateCommand {
			return a.createPhoneNumber(ctx, c)
		} else {
			return a.updatePhoneNumber(ctx, c)
		}
	default:
		return errors.New("invalid command type")
	}
}

func (a *PhoneNumberAggregate) createPhoneNumber(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.createPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	event, err := events.NewPhoneNumberCreateEvent(a, cmd.Tenant, cmd.RawPhoneNumber, cmd.Source, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberCreateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) updatePhoneNumber(ctx context.Context, cmd *command.UpsertPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.updatePhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.LogFields(log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())
	if cmd.Source.SourceOfTruth == "" {
		cmd.Source.SourceOfTruth = a.PhoneNumber.Source.SourceOfTruth
	}

	event, err := events.NewPhoneNumberUpdateEvent(a, cmd.Tenant, cmd.Source.SourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberUpdateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.UserID)

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) FailedPhoneNumberValidation(ctx context.Context, tenant, rawPhoneNumber, countryCodeA2, validationError string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.FailedPhoneNumberValidation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewPhoneNumberFailedValidationEvent(a, tenant, rawPhoneNumber, countryCodeA2, validationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberFailedValidationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) SkippedPhoneNumberValidation(ctx context.Context, tenant, rawPhoneNumber, countryCodeA2, validationSkipReason string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.SkippedPhoneNumberValidation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewPhoneNumberSkippedValidationEvent(a, tenant, rawPhoneNumber, countryCodeA2, validationSkipReason)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberSkippedValidationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) PhoneNumberValidated(ctx context.Context, tenant, rawPhoneNumber, e164, countryCodeA2 string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.PhoneNumberValidated")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewPhoneNumberValidatedEvent(a, tenant, rawPhoneNumber, e164, countryCodeA2)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberValidatedEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, "")

	return a.Apply(event)
}
