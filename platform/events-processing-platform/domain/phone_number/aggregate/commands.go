package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/domain/phone_number/events"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *PhoneNumberAggregate) CreatePhoneNumber(ctx context.Context, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.CreatePhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTime(createdAt, func() time.Time { return utils.Now() })
	updatedAtNotNil := utils.IfNotNilTime(updatedAt, func() time.Time { return createdAtNotNil })
	event, err := events.NewPhoneNumberCreatedEvent(a, tenant, rawPhoneNumber, source, sourceOfTruth, appSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberCreatedEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}

func (a *PhoneNumberAggregate) UpdatePhoneNumber(ctx context.Context, tenant, sourceOfTruth string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "PhoneNumberAggregate.UpdatePhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTime(updatedAt, func() time.Time { return utils.Now() })

	event, err := events.NewPhoneNumberUpdatedEvent(a, tenant, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewPhoneNumberUpdatedEvent")
	}

	// FIXME alexb check what type of metadata should be set into event and apply it to all aggregation commands
	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
