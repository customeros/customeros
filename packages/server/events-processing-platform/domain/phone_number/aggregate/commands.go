package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *PhoneNumberAggregate) CreatePhoneNumber(ctx context.Context, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.CreatePhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(createdAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, createdAtNotNil)
	event, err := events.NewPhoneNumberCreateEvent(a, tenant, rawPhoneNumber, source, sourceOfTruth, appSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) UpdatePhoneNumber(ctx context.Context, tenant, sourceOfTruth string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.UpdatePhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())
	if sourceOfTruth == "" {
		sourceOfTruth = a.PhoneNumber.Source.SourceOfTruth
	}

	event, err := events.NewPhoneNumberUpdateEvent(a, tenant, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

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

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

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

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

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

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
