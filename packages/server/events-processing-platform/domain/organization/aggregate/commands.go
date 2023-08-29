package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	locerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

func (a *OrganizationAggregate) HandleCommand(ctx context.Context, command eventstore.Command) error {
	switch c := command.(type) {
	case *cmd.RequestNextCycleDateCommand:
		return a.requestNextCycleDate(ctx, c)
	case *cmd.RequestRenewalForecastCommand:
		return a.requestRenewalForecast(ctx, c)
	case *cmd.UpdateRenewalLikelihoodCommand:
		return a.updateRenewalLikelihood(ctx, c)
	case *cmd.UpdateRenewalForecastCommand:
		return a.updateRenewalForecast(ctx, c)
	case *cmd.UpdateBillingDetailsCommand:
		return a.updateBillingDetails(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *OrganizationAggregate) CreateOrganization(ctx context.Context, organizationFields *models.OrganizationFields) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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
		return locerr.ErrPhoneNumberNotFound
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
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.Now()

	email, ok := a.Organization.Emails[emailId]
	if !ok {
		return locerr.ErrEmailNotFound
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
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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

func (a *OrganizationAggregate) AddSocial(ctx context.Context, tenant, socialId, platformName, url string, source common_models.Source, createdAt *time.Time, updatedAt *time.Time) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.AddSocial")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

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

func (a *OrganizationAggregate) updateRenewalLikelihood(ctx context.Context, command *cmd.UpdateRenewalLikelihoodCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.updateRenewalLikelihood")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAt := command.Fields.UpdatedAt
	if updatedAt == utils.ZeroTime() {
		updatedAt = utils.Now()
	}

	event, err := events.NewOrganizationUpdateRenewalLikelihoodEvent(a, command.Fields.RenewalLikelihood, a.Organization.RenewalLikelihood.RenewalLikelihood, command.Fields.UpdatedBy, command.Fields.Comment, updatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = a.Tenant
	metadata["user-id"] = command.Fields.UpdatedBy
	if err = event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) updateRenewalForecast(ctx context.Context, command *cmd.UpdateRenewalForecastCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.updateRenewalForecast")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAt := command.Fields.UpdatedAt
	if updatedAt == utils.ZeroTime() {
		updatedAt = utils.Now()
	}

	event, err := events.NewOrganizationUpdateRenewalForecastEvent(a, command.Fields.Amount, command.Fields.PotentialAmount, command.Fields.UpdatedBy, command.Fields.Comment, updatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = a.Tenant
	if command.Fields.UpdatedBy != "" {
		metadata["user-id"] = command.Fields.UpdatedBy
	}
	if err = event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) requestRenewalForecast(ctx context.Context, command *cmd.RequestRenewalForecastCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.RequestRenewalForecast")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationRequestRenewalForecastEvent(a, a.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestRenewalForecastEvent")
	}

	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = a.Tenant
	if err = event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) updateBillingDetails(ctx context.Context, command *cmd.UpdateBillingDetailsCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.updateBillingDetails")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationUpdateBillingDetailsEvent(a, command.Fields.Amount, command.Fields.Frequency, command.Fields.RenewalCycle, command.Fields.UpdatedBy, command.Fields.RenewalCycleStart, command.Fields.RenewalCycleNext)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = a.Tenant
	if command.Fields.UpdatedBy != "" {
		metadata["user-id"] = command.Fields.UpdatedBy
	}
	if err = event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}

func (a *OrganizationAggregate) requestNextCycleDate(ctx context.Context, command *cmd.RequestNextCycleDateCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.requestNextCycleDate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationRequestNextCycleDateEvent(a, a.Tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestNextCycleDateEvent")
	}

	metadata := tracing.ExtractTextMapCarrier(span.Context())
	metadata["tenant"] = a.Tenant
	if err = event.SetMetadata(metadata); err != nil {
		tracing.TraceErr(span, err)
	}

	return a.Apply(event)
}
