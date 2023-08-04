package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	local_errors "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *OrganizationAggregate) CreateOrganization(ctx context.Context, organizationFields *models.OrganizationFields) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", organizationFields.Tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, createdAtNotNil)
	event, err := events.NewOrganizationCreateEvent(a, organizationFields, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationCreateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) UpdateOrganization(ctx context.Context, organizationFields *models.OrganizationFields) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", organizationFields.Tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, utils.Now())
	if organizationFields.Source.SourceOfTruth == "" {
		organizationFields.Source.SourceOfTruth = a.Organization.Source.SourceOfTruth
	}

	event, err := events.NewOrganizationUpdateEvent(a, organizationFields, updatedAtNotNil, organizationFields.IgnoreEmptyFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUpdateEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) LinkPhoneNumber(ctx context.Context, tenant, phoneNumberId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkPhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkPhoneNumberEvent(a, tenant, phoneNumberId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Organization.PhoneNumbers[phoneNumberId]
	if !ok {
		return local_errors.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := events.NewOrganizationLinkPhoneNumberEvent(a, tenant, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}

func (a *OrganizationAggregate) LinkEmail(ctx context.Context, tenant, emailId, label string, primary bool) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkEmailEvent(a, tenant, emailId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkEmailEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	email, ok := a.Organization.Emails[emailId]
	if !ok {
		return local_errors.ErrEmailNotFound
	}

	if email.Primary {
		event, err := events.NewOrganizationLinkEmailEvent(a, tenant, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationLinkEmailEvent")
		}

		if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
			tracing.TraceErr(span, err)
		}
		return a.Apply(event)
	}
	return nil
}

func (a *OrganizationAggregate) LinkDomain(ctx context.Context, tenant, domain string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkDomain")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	event, err := events.NewOrganizationLinkDomainEvent(a, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkDomainEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) AddSocial(ctx context.Context, tenant, socialId, platformName, url string, source commonModels.Source, createdAt *time.Time, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.AddSocial")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(createdAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(updatedAt, createdAtNotNil)
	localSource := utils.StringFirstNonEmpty(source.Source, constants.SourceOpenline)
	localSourceOfTruth := utils.StringFirstNonEmpty(source.SourceOfTruth, constants.SourceOpenline)
	localAppSource := utils.StringFirstNonEmpty(source.AppSource, constants.AppSourceEventProcessingPlatform)

	event, err := events.NewOrganizationAddSocialEvent(a, tenant, socialId, platformName, url, localSource, localSourceOfTruth, localAppSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}

	if err = event.SetMetadata(tracing.ExtractTextMapCarrier(span.Context())); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
