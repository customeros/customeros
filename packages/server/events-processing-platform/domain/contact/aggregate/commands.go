package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	localErrors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *ContactAggregate) CreateContact(ctx context.Context, contactDto *models.ContactDto) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.CreateContact")
	defer span.Finish()
	span.LogFields(log.String("Tenant", contactDto.Tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(contactDto.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(contactDto.UpdatedAt, createdAtNotNil)
	event, err := events.NewContactCreateEvent(a, contactDto, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *ContactAggregate) UpdateContact(ctx context.Context, contactDto *models.ContactDto) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.UpdateContact")
	defer span.Finish()
	span.LogFields(log.String("Tenant", contactDto.Tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(contactDto.UpdatedAt, utils.Now())
	if contactDto.Source.SourceOfTruth == "" {
		contactDto.Source.SourceOfTruth = a.Contact.Source.SourceOfTruth
	}

	event, err := events.NewContactUpdateEvent(a, contactDto, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *ContactAggregate) LinkPhoneNumber(ctx context.Context, tenant, phoneNumberId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.LinkPhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewContactLinkPhoneNumberEvent(a, tenant, phoneNumberId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkPhoneNumberEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *ContactAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Contact.PhoneNumbers[phoneNumberId]
	if !ok {
		return localErrors.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := events.NewContactLinkPhoneNumberEvent(a, tenant, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewContactLinkPhoneNumberEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}

func (a *ContactAggregate) LinkEmail(ctx context.Context, tenant, emailId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.LinkEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewContactLinkEmailEvent(a, tenant, emailId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactLinkEmailEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *ContactAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	email, ok := a.Contact.Emails[emailId]
	if !ok {
		return localErrors.ErrEmailNotFound
	}

	if email.Primary {
		event, err := events.NewContactLinkEmailEvent(a, tenant, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewContactLinkEmailEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}
