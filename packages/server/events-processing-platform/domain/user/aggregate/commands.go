package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	local_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func (a *UserAggregate) CreateUser(ctx context.Context, userDto *models.UserFields) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.CreateUser")
	defer span.Finish()
	span.LogFields(log.String("Tenant", userDto.Tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(userDto.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(userDto.UpdatedAt, createdAtNotNil)
	event, err := events.NewUserCreateEvent(a, userDto, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *UserAggregate) UpdateUser(ctx context.Context, userDto *models.UserFields) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.UpdateUser")
	defer span.Finish()
	span.LogFields(log.String("Tenant", userDto.Tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(userDto.UpdatedAt, utils.Now())
	if userDto.Source.SourceOfTruth == "" {
		userDto.Source.SourceOfTruth = a.User.Source.SourceOfTruth
	}

	event, err := events.NewUserUpdateEvent(a, userDto, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *UserAggregate) LinkJobRole(ctx context.Context, tenant, jobRoleId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkJobRole")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkJobRoleEvent(a, tenant, jobRoleId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkJobRoleEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *UserAggregate) LinkPhoneNumber(ctx context.Context, tenant, phoneNumberId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkPhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkPhoneNumberEvent(a, tenant, phoneNumberId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *UserAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.User.PhoneNumbers[phoneNumberId]
	if !ok {
		return local_errors.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := events.NewUserLinkPhoneNumberEvent(a, tenant, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewUserLinkPhoneNumberEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}

func (a *UserAggregate) LinkEmail(ctx context.Context, tenant, emailId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.LinkEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewUserLinkEmailEvent(a, tenant, emailId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUserLinkEmailEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *UserAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	email, ok := a.User.Emails[emailId]
	if !ok {
		return local_errors.ErrEmailNotFound
	}

	if email.Primary {
		event, err := events.NewUserLinkEmailEvent(a, tenant, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewUserLinkEmailEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}
