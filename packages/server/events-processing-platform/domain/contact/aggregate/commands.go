package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *ContactAggregate) CreateContact(ctx context.Context, tenant, firstName, lastName, name, prefix, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.CreateContact")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(createdAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, createdAtNotNil)
	event, err := events.NewContactCreatedEvent(a, tenant, firstName, lastName, name, prefix, source, sourceOfTruth, appSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactCreatedEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}

func (a *ContactAggregate) UpdateContact(ctx context.Context, tenant, sourceOfTruth, firstName, lastName, name, prefix string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.UpdateContact")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())
	if sourceOfTruth == "" {
		sourceOfTruth = a.Contact.Source.SourceOfTruth
	}

	event, err := events.NewContactUpdatedEvent(a, tenant, firstName, lastName, name, prefix, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactUpdatedEvent")
	}

	// FIXME alexb check what type of metadata should be set into event and apply it to all aggregation commands
	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "SetMetadata")
	}

	return a.Apply(event)
}
