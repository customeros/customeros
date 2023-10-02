package aggregate

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	locerr "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	case *cmd.LinkDomainCommand:
		return a.linkDomain(ctx, c)
	case *cmd.AddSocialCommand:
		return a.addSocial(ctx, c)
	case *cmd.HideOrganizationCommand:
		return a.hideOrganization(ctx, c)
	case *cmd.ShowOrganizationCommand:
		return a.showOrganization(ctx, c)
	case *cmd.RefreshLastTouchpointCommand:
		return a.refreshLastTouchpoint(ctx, c)
	case *cmd.UpsertCustomFieldCommand:
		return a.upsertCustomField(ctx, c)
	default:
		return errors.New("invalid command type")
	}
}

func (a *OrganizationAggregate) CreateOrganization(ctx context.Context, organizationFields *models.OrganizationFields, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	var eventsOnCreate []eventstore.Event

	createdAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, createdAtNotNil)

	createEvent, err := events.NewOrganizationCreateEvent(a, organizationFields, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationCreateEvent")
	}
	aggregate.EnrichEventWithMetadata(&createEvent, &span, a.Tenant, userId)
	eventsOnCreate = append(eventsOnCreate, createEvent)

	if organizationFields.OrganizationDataFields.Website != "" {
		webscrapeEvent, err := events.NewOrganizationRequestScrapeByWebsite(a, organizationFields.OrganizationDataFields.Website)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationCreateEvent")
		}
		aggregate.EnrichEventWithMetadata(&webscrapeEvent, &span, a.Tenant, userId)
		eventsOnCreate = append(eventsOnCreate, webscrapeEvent)
	}

	return a.ApplyAll(eventsOnCreate)
}

func (a *OrganizationAggregate) UpdateOrganization(ctx context.Context, organizationFields *models.OrganizationFields, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	var eventsOnUpdate []eventstore.Event

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, utils.Now())
	if organizationFields.Source.SourceOfTruth == "" {
		organizationFields.Source.SourceOfTruth = a.Organization.Source.SourceOfTruth
	}

	event, err := events.NewOrganizationUpdateEvent(a, organizationFields, updatedAtNotNil, organizationFields.IgnoreEmptyFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUpdateEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)
	eventsOnUpdate = append(eventsOnUpdate, event)

	// if website updated, request webscrape
	if organizationFields.OrganizationDataFields.Website != "" && organizationFields.OrganizationDataFields.Website != a.Organization.Website {
		webscrapeEvent, err := events.NewOrganizationRequestScrapeByWebsite(a, organizationFields.OrganizationDataFields.Website)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationCreateEvent")
		}
		aggregate.EnrichEventWithMetadata(&webscrapeEvent, &span, a.Tenant, userId)
		eventsOnUpdate = append(eventsOnUpdate, webscrapeEvent)
	}

	return a.ApplyAll(eventsOnUpdate)
}

func (a *OrganizationAggregate) LinkPhoneNumber(ctx context.Context, tenant, phoneNumberId, label string, primary bool, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkPhoneNumber")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkPhoneNumberEvent(a, phoneNumberId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("AggregateID", a.GetID()))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Organization.PhoneNumbers[phoneNumberId]
	if !ok {
		return locerr.ErrPhoneNumberNotFound
	}

	if phoneNumber.Primary {
		event, err := events.NewOrganizationLinkPhoneNumberEvent(a, phoneNumberId, phoneNumber.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
		}

		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)
		return a.Apply(event)
	}
	return nil
}

func (a *OrganizationAggregate) LinkEmail(ctx context.Context, tenant, emailId, label string, primary bool, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkEmail")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkEmailEvent(a, emailId, label, primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkEmailEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) SetEmailNonPrimary(ctx context.Context, tenant, emailId, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.Now()

	email, ok := a.Organization.Emails[emailId]
	if !ok {
		return locerr.ErrEmailNotFound
	}

	if email.Primary {
		event, err := events.NewOrganizationLinkEmailEvent(a, emailId, email.Label, false, updatedAtNotNil)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationLinkEmailEvent")
		}

		aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, userId)
		return a.Apply(event)
	}
	return nil
}

func (a *OrganizationAggregate) linkDomain(ctx context.Context, command *cmd.LinkDomainCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.linkDomain")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationLinkDomainEvent(a, command.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkDomainEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) addSocial(ctx context.Context, command *cmd.AddSocialCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addSocial")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(command.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, createdAtNotNil)
	localSource := utils.StringFirstNonEmpty(command.Source.Source, constants.SourceOpenline)
	localSourceOfTruth := utils.StringFirstNonEmpty(command.Source.SourceOfTruth, constants.SourceOpenline)
	localAppSource := utils.StringFirstNonEmpty(command.Source.AppSource, constants.AppSourceEventProcessingPlatform)

	event, err := events.NewOrganizationAddSocialEvent(a, command.SocialId, command.SocialPlatform, command.SocialUrl, localSource, localSourceOfTruth, localAppSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

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

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

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

	event, err := events.NewOrganizationUpdateRenewalForecastEvent(a, command.Fields.Amount, command.Fields.PotentialAmount, a.Organization.RenewalForecast.Amount, command.Fields.UpdatedBy, command.Fields.Comment, updatedAt, a.Organization.RenewalLikelihood.RenewalLikelihood)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) requestRenewalForecast(ctx context.Context, command *cmd.RequestRenewalForecastCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.requestRenewalForecast")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationRequestRenewalForecastEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestRenewalForecastEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

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

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) requestNextCycleDate(ctx context.Context, command *cmd.RequestNextCycleDateCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.requestNextCycleDate")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationRequestNextCycleDateEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestNextCycleDateEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) hideOrganization(ctx context.Context, command *cmd.HideOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.hideOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewHideOrganizationEventEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewHideOrganizationEventEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) showOrganization(ctx context.Context, command *cmd.ShowOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.showOrganization")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewShowOrganizationEventEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewShowOrganizationEventEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) refreshLastTouchpoint(ctx context.Context, command *cmd.RefreshLastTouchpointCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshLastTouchpoint")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewOrganizationRefreshLastTouchpointEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshLastTouchpointEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, command.UserID)

	return a.Apply(event)
}

func (a *OrganizationAggregate) upsertCustomField(ctx context.Context, command *cmd.UpsertCustomFieldCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.upsertCustomField")
	defer span.Finish()
	span.LogFields(log.String("Tenant", a.Tenant), log.String("AggregateID", a.GetID()), log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(command.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(command.UpdatedAt, createdAtNotNil)
	sourceFields := command.Source
	if sourceFields.Source == "" {
		sourceFields.Source = constants.SourceOpenline
	}
	if sourceFields.SourceOfTruth == "" {
		if val, ok := a.Organization.CustomFields[command.CustomFieldData.Id]; ok {
			sourceFields.SourceOfTruth = val.Source.SourceOfTruth
		} else {
			sourceFields.SourceOfTruth = constants.SourceOpenline
		}
	}
	if sourceFields.AppSource == "" {
		sourceFields.AppSource = constants.AppSourceEventProcessingPlatform
	}

	found := false
	if _, ok := a.Organization.CustomFields[command.CustomFieldData.Id]; ok {
		found = true
	}

	event, err := events.NewOrganizationUpsertCustomField(a, sourceFields, createdAtNotNil, updatedAtNotNil, command.CustomFieldData, found)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUpsertCustomField")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, command.Tenant, command.UserID)

	return a.Apply(event)
}
