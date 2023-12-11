package aggregate

import (
	"context"
	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	localerror "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

func (a *OrganizationAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.LinkDomainCommand:
		return a.linkDomain(ctx, c)
	case *command.AddSocialCommand:
		return a.addSocial(ctx, c)
	case *command.HideOrganizationCommand:
		return a.hideOrganization(ctx, c)
	case *command.ShowOrganizationCommand:
		return a.showOrganization(ctx, c)
	case *command.RefreshLastTouchpointCommand:
		return a.refreshLastTouchpoint(ctx, c)
	case *command.RefreshArrCommand:
		return a.refreshArr(ctx, c)
	case *command.RefreshRenewalSummaryCommand:
		return a.refreshRenewalSummary(ctx, c)
	case *command.UpsertCustomFieldCommand:
		return a.upsertCustomField(ctx, c)
	case *command.LinkEmailCommand:
		return a.linkEmail(ctx, c)
	case *command.LinkPhoneNumberCommand:
		return a.linkPhoneNumber(ctx, c)
	case *command.LinkLocationCommand:
		return a.linkLocation(ctx, c)
	case *command.AddParentCommand:
		return a.addParentOrganization(ctx, c)
	case *command.RemoveParentCommand:
		return a.removeParentOrganization(ctx, c)
	case *command.WebScrapeOrganizationCommand:
		return a.webScrapeOrganization(ctx, c)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *OrganizationAggregate) CreateOrganization(ctx context.Context, organizationFields *model.OrganizationFields, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "organizationFields", organizationFields)

	var eventsOnCreate []eventstore.Event

	createdAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, createdAtNotNil)
	organizationFields.Source.SetDefaultValues()

	createEvent, err := events.NewOrganizationCreateEvent(a, organizationFields, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationCreateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&createEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: userId,
		App:    organizationFields.Source.AppSource,
	})
	eventsOnCreate = append(eventsOnCreate, createEvent)

	if organizationFields.OrganizationDataFields.Website != "" && strings.Contains(organizationFields.OrganizationDataFields.Website, ".") {
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

func (a *OrganizationAggregate) UpdateOrganization(ctx context.Context, organizationFields *model.OrganizationFields, loggedInUserId, webScrapedUrl string, fieldsMask []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()), log.String("loggedInUserId", loggedInUserId), log.Object("fieldsMask", fieldsMask))
	tracing.LogObjectAsJson(span, "organizationFields", organizationFields)

	var eventsOnUpdate []eventstore.Event

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, utils.Now())

	event, err := events.NewOrganizationUpdateEvent(a, organizationFields, updatedAtNotNil, webScrapedUrl, fieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUpdateEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: loggedInUserId,
		App:    organizationFields.Source.AppSource,
	})
	eventsOnUpdate = append(eventsOnUpdate, event)

	// if website updated, request web scrape by website
	websiteChanged := organizationFields.OrganizationDataFields.Website != a.Organization.Website && (len(fieldsMask) == 0 || utils.Contains(fieldsMask, model.FieldMaskWebsite))
	if organizationFields.OrganizationDataFields.Website != "" && websiteChanged && strings.Contains(organizationFields.OrganizationDataFields.Website, ".") {
		webScrapeEvent, err := events.NewOrganizationRequestScrapeByWebsite(a, organizationFields.OrganizationDataFields.Website)
		if err != nil {
			tracing.TraceErr(span, err)
			return errors.Wrap(err, "NewOrganizationCreateEvent")
		}
		aggregate.EnrichEventWithMetadataExtended(&webScrapeEvent, span, aggregate.EventMetadata{
			Tenant: a.Tenant,
			UserId: loggedInUserId,
			App:    organizationFields.Source.AppSource,
		})
		eventsOnUpdate = append(eventsOnUpdate, webScrapeEvent)
	}

	return a.ApplyAll(eventsOnUpdate)
}

func (a *OrganizationAggregate) linkPhoneNumber(ctx context.Context, cmd *command.LinkPhoneNumberCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.linkPhoneNumber")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkPhoneNumberEvent(a, cmd.PhoneNumberId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkPhoneNumberEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.Organization.PhoneNumbers {
			if k != cmd.PhoneNumberId && v.Primary {
				if err = a.SetPhoneNumberNonPrimary(ctx, cmd.Tenant, k, cmd.LoggedInUserId); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *OrganizationAggregate) SetPhoneNumberNonPrimary(ctx context.Context, tenant, phoneNumberId, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetPhoneNumberNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("phoneNumberId", phoneNumberId), log.String("userId", userId))

	updatedAtNotNil := utils.Now()

	phoneNumber, ok := a.Organization.PhoneNumbers[phoneNumberId]
	if !ok {
		return localerror.ErrPhoneNumberNotFound
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

func (a *OrganizationAggregate) linkEmail(ctx context.Context, cmd *command.LinkEmailCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.linkEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkEmailEvent(a, cmd.EmailId, cmd.Label, cmd.Primary, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkEmailEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	err = a.Apply(event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if cmd.Primary {
		for k, v := range a.Organization.Emails {
			if k != cmd.EmailId && v.Primary {
				if err = a.SetEmailNonPrimary(ctx, k, cmd.LoggedInUserId); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (a *OrganizationAggregate) linkLocation(ctx context.Context, cmd *command.LinkLocationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.linkLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.Now()

	event, err := events.NewOrganizationLinkLocationEvent(a, cmd.LocationId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkLocationEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) SetEmailNonPrimary(ctx context.Context, emailId, userId string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.SetEmailNonPrimary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()), log.String("emailId", emailId), log.String("userId", userId))

	updatedAtNotNil := utils.Now()

	email, ok := a.Organization.Emails[emailId]
	if !ok {
		return localerror.ErrEmailNotFound
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

func (a *OrganizationAggregate) linkDomain(ctx context.Context, cmd *command.LinkDomainCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.linkDomain")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	if aggregate.AllowCheckIfEventIsRedundant(cmd.AppSource, cmd.LoggedInUserId) {
		if utils.Contains(a.Organization.Domains, strings.TrimSpace(cmd.Domain)) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	event, err := events.NewOrganizationLinkDomainEvent(a, cmd.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationLinkDomainEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) addSocial(ctx context.Context, cmd *command.AddSocialCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addSocial")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	localSource := utils.StringFirstNonEmpty(cmd.Source.Source, constants.SourceOpenline)
	localSourceOfTruth := utils.StringFirstNonEmpty(cmd.Source.SourceOfTruth, constants.SourceOpenline)
	localAppSource := utils.StringFirstNonEmpty(cmd.Source.AppSource, constants.AppSourceEventProcessingPlatform)

	if existingSocialId := a.Organization.GetSocialIdForUrl(cmd.SocialUrl); existingSocialId != "" {
		cmd.SocialId = existingSocialId
	} else if cmd.SocialId == "" {
		cmd.SocialId = uuid.New().String()
	}

	event, err := events.NewOrganizationAddSocialEvent(a, cmd.SocialId, cmd.SocialPlatform, cmd.SocialUrl, localSource, localSourceOfTruth, localAppSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) hideOrganization(ctx context.Context, cmd *command.HideOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.hideOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewHideOrganizationEventEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewHideOrganizationEventEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) showOrganization(ctx context.Context, cmd *command.ShowOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.showOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewShowOrganizationEventEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewShowOrganizationEventEvent")
	}

	aggregate.EnrichEventWithMetadata(&event, &span, a.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) refreshLastTouchpoint(ctx context.Context, cmd *command.RefreshLastTouchpointCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshLastTouchpoint")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationRefreshLastTouchpointEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshLastTouchpointEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) refreshArr(ctx context.Context, cmd *command.RefreshArrCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationRefreshArrEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshArrEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) refreshRenewalSummary(ctx context.Context, cmd *command.RefreshRenewalSummaryCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshRenewalSummary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationRefreshRenewalSummaryEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshArrEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) upsertCustomField(ctx context.Context, cmd *command.UpsertCustomFieldCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.upsertCustomField")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(cmd.CreatedAt, utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, createdAtNotNil)
	sourceFields := cmd.Source
	if sourceFields.Source == "" {
		sourceFields.Source = constants.SourceOpenline
	}
	if sourceFields.SourceOfTruth == "" {
		if val, ok := a.Organization.CustomFields[cmd.CustomFieldData.Id]; ok {
			sourceFields.SourceOfTruth = val.Source.SourceOfTruth
		} else {
			sourceFields.SourceOfTruth = constants.SourceOpenline
		}
	}
	if sourceFields.AppSource == "" {
		sourceFields.AppSource = constants.AppSourceEventProcessingPlatform
	}

	found := false
	if _, ok := a.Organization.CustomFields[cmd.CustomFieldData.Id]; ok {
		found = true
	}

	event, err := events.NewOrganizationUpsertCustomField(a, sourceFields, createdAtNotNil, updatedAtNotNil, cmd.CustomFieldData, found)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUpsertCustomField")
	}
	aggregate.EnrichEventWithMetadata(&event, &span, cmd.Tenant, cmd.LoggedInUserId)

	return a.Apply(event)
}

func (a *OrganizationAggregate) addParentOrganization(ctx context.Context, cmd *command.AddParentCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addParentOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationAddParentEvent(a, cmd.ParentOrganizationId, cmd.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddParentEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: cmd.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) removeParentOrganization(ctx context.Context, cmd *command.RemoveParentCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.removeParentOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationRemoveParentEvent(a, cmd.ParentOrganizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRemoveParentEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: cmd.Tenant,
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) webScrapeOrganization(ctx context.Context, cmd *command.WebScrapeOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.webScrapeOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	event, err := events.NewOrganizationRequestScrapeByWebsite(a, cmd.Website)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRequestScrapeByWebsite")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}
