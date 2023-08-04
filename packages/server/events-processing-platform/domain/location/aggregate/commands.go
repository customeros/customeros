package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	models_common "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *LocationAggregate) CreateLocation(ctx context.Context, tenant, name, rawAddress string, locationAddress models.LocationAddress, source models_common.Source, createdAt, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.CreateLocation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(createdAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, createdAtNotNil)
	event, err := events.NewLocationCreateEvent(a, tenant, name, rawAddress, source.Source, source.SourceOfTruth, source.AppSource, createdAtNotNil, updatedAtNotNil, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *LocationAggregate) UpdateLocation(ctx context.Context, tenant, name, rawAddress string, locationAddress models.LocationAddress, sourceOfTruth string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.UpdateLocation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())
	if sourceOfTruth == "" {
		sourceOfTruth = a.Location.Source.SourceOfTruth
	}

	event, err := events.NewLocationUpdateEvent(a, tenant, name, rawAddress, sourceOfTruth, updatedAtNotNil, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *LocationAggregate) FailLocationValidation(ctx context.Context, tenant, rawAddress, country, validationError string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.FailLocationValidation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewLocationFailedValidationEvent(a, tenant, rawAddress, country, validationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationFailedValidationEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *LocationAggregate) SkipLocationValidation(ctx context.Context, tenant, rawAddress, validationSkipReason string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.SkipLocationValidation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewLocationSkippedValidationEvent(a, tenant, rawAddress, validationSkipReason)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationSkippedValidationEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *LocationAggregate) LocationValidated(ctx context.Context, tenant, rawAddress, countryForValidation string, locationAddress models.LocationAddress) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "LocationAggregate.LocationValidated")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewLocationValidatedEvent(a, tenant, rawAddress, countryForValidation, locationAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLocationValidatedEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
