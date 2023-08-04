package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *EmailAggregate) CreateEmail(ctx context.Context, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.CreateEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(createdAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, createdAtNotNil)
	event, err := events.NewEmailCreateEvent(a, tenant, rawEmail, source, sourceOfTruth, appSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *EmailAggregate) UpdateEmail(ctx context.Context, rawEmail, tenant, sourceOfTruth string, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.UpdateEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, utils.Now())
	if sourceOfTruth == "" {
		sourceOfTruth = a.Email.Source.SourceOfTruth
	}

	event, err := events.NewEmailUpdateEvent(a, rawEmail, tenant, sourceOfTruth, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *EmailAggregate) FailEmailValidation(ctx context.Context, tenant, validationError string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.FailEmailValidation")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewEmailFailedValidationEvent(a, tenant, validationError)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailFailedValidationEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *EmailAggregate) EmailValidated(ctx context.Context, tenant, rawEmail, isReachable, validationError, domain, username, emailAddress string,
	acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, IsDeliverable, isDisabled, isValidSyntax bool) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "EmailAggregate.EmailValidated")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewEmailValidatedEvent(a, tenant, rawEmail, isReachable, validationError, domain, username, emailAddress, acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, IsDeliverable, isDisabled, isValidSyntax)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewEmailValidatedEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
