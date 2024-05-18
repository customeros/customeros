package aggregate

import (
	"context"
	"strings"

	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"

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
)

func (a *OrganizationAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.CreateBillingProfileGrpcRequest:
		return a.CreateBillingProfile(ctx, r)
	case *organizationpb.UpdateBillingProfileGrpcRequest:
		return nil, a.UpdateBillingProfile(ctx, r)
	case *organizationpb.LinkEmailToBillingProfileGrpcRequest:
		return nil, a.LinkEmailToBillingProfile(ctx, r)
	case *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest:
		return nil, a.UnlinkEmailFromBillingProfile(ctx, r)
	case *organizationpb.LinkLocationToBillingProfileGrpcRequest:
		return nil, a.LinkLocationToBillingProfile(ctx, r)
	case *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest:
		return nil, a.UnlinkLocationFromBillingProfile(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationTempAggregate) HandleRequest(ctx context.Context, request any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.RefreshRenewalSummaryGrpcRequest:
		return nil, a.refreshRenewalSummary(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationAggregate) CreateBillingProfile(ctx context.Context, request *organizationpb.CreateBillingProfileGrpcRequest) (billingProfileId string, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.CreateBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())
	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), createdAtNotNil)
	sourceFields := commonmodel.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	billingProfileId = utils.NewUUIDIfEmpty(request.BillingProfileId)

	event, err := events.NewBillingProfileCreateEvent(a, billingProfileId, request.LegalName, request.TaxId, sourceFields, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewBillingProfileCreateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return billingProfileId, a.Apply(event)
}

func (a *OrganizationAggregate) UpdateBillingProfile(ctx context.Context, request *organizationpb.UpdateBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.UpdatedAt), utils.Now())
	var fieldsMask []string
	if utils.ContainsElement(request.FieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_LEGAL_NAME) {
		fieldsMask = append(fieldsMask, events.FieldMaskLegalName)
	}
	if utils.ContainsElement(request.FieldsMask, organizationpb.BillingProfileFieldMask_BILLING_PROFILE_PROPERTY_TAX_ID) {
		fieldsMask = append(fieldsMask, events.FieldMaskTaxId)
	}

	updateEvent, err := events.NewBillingProfileUpdateEvent(a, request.BillingProfileId, request.LegalName, request.TaxId, updatedAtNotNil, fieldsMask)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewBillingProfileUpdateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&updateEvent, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(updateEvent)
}

func (a *OrganizationAggregate) LinkEmailToBillingProfile(ctx context.Context, request *organizationpb.LinkEmailToBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkEmailToBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewLinkEmailToBillingProfileEvent(a, request.BillingProfileId, request.EmailId, request.Primary, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLinkEmailToBillingProfileEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) UnlinkEmailFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkEmailFromBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UnlinkEmailFromBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewUnlinkEmailFromBillingProfileEvent(a, request.BillingProfileId, request.EmailId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUnlinkEmailFromBillingProfileEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) LinkLocationToBillingProfile(ctx context.Context, request *organizationpb.LinkLocationToBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.LinkLocationToBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewLinkLocationToBillingProfileEvent(a, request.BillingProfileId, request.LocationId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewLinkLocationToBillingProfileEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) UnlinkLocationFromBillingProfile(ctx context.Context, request *organizationpb.UnlinkLocationFromBillingProfileGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UnlinkLocationFromBillingProfile")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()))

	event, err := events.NewUnlinkLocationFromBillingProfileEvent(a, request.BillingProfileId, request.LocationId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUnlinkLocationFromBillingProfileEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: request.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

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
	case *command.UpdateOnboardingStatusCommand:
		return a.updateOnboardingStatus(ctx, c)
	case *command.UpdateOrganizationOwnerCommand:
		return a.UpdateOrganizationOwner(ctx, c)

	default:
		tracing.TraceErr(span, eventstore.ErrInvalidCommandType)
		return eventstore.ErrInvalidCommandType
	}
}

func (a *OrganizationTempAggregate) HandleCommand(ctx context.Context, cmd eventstore.Command) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.HandleCommand")
	defer span.Finish()

	switch c := cmd.(type) {
	case *command.RefreshLastTouchpointCommand:
		return a.refreshLastTouchpoint(ctx, c)
	case *command.RefreshArrCommand:
		return a.refreshArr(ctx, c)

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

func (a *OrganizationAggregate) UpdateOrganization(ctx context.Context, organizationFields *model.OrganizationFields, loggedInUserId, enrichDomain, enrichSource string, fieldsMask []string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("AggregateVersion", a.GetVersion()), log.String("loggedInUserId", loggedInUserId), log.Object("fieldsMask", fieldsMask))
	tracing.LogObjectAsJson(span, "organizationFields", organizationFields)

	if aggregate.AllowCheckForNoChanges(organizationFields.Source.AppSource, loggedInUserId) {
		if a.Organization.SkipUpdate(organizationFields) {
			span.SetTag(tracing.SpanTagRedundantEventSkipped, true)
			return nil
		}
	}

	var eventsOnUpdate []eventstore.Event

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(organizationFields.UpdatedAt, utils.Now())

	event, err := events.NewOrganizationUpdateEvent(a, organizationFields, updatedAtNotNil, enrichDomain, enrichSource, fieldsMask)
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

	if aggregate.AllowCheckForNoChanges(cmd.AppSource, cmd.LoggedInUserId) {
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
	} else {
		cmd.SocialId = utils.NewUUIDIfEmpty(cmd.SocialId)
	}

	event, err := events.NewOrganizationAddSocialEvent(a, cmd.SocialId, cmd.SocialUrl, localSource, localSourceOfTruth, localAppSource, createdAtNotNil, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.Source.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) hideOrganization(ctx context.Context, cmd *command.HideOrganizationCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.hideOrganization")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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

func (a *OrganizationTempAggregate) refreshLastTouchpoint(ctx context.Context, cmd *command.RefreshLastTouchpointCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationTempAggregate.refreshLastTouchpoint")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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

func (a *OrganizationTempAggregate) refreshArr(ctx context.Context, cmd *command.RefreshArrCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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

func (a *OrganizationTempAggregate) refreshRenewalSummary(ctx context.Context, request *organizationpb.RefreshRenewalSummaryGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.refreshRenewalSummary")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	event, err := events.NewOrganizationRefreshRenewalSummaryEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRefreshRenewalSummaryEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) upsertCustomField(ctx context.Context, cmd *command.UpsertCustomFieldCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.upsertCustomField")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
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

func (a *OrganizationAggregate) updateOnboardingStatus(ctx context.Context, cmd *command.UpdateOnboardingStatusCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.updateOnboardingStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAtNotNil := utils.IfNotNilTimeWithDefault(cmd.UpdatedAt, utils.Now())

	event, err := events.NewUpdateOnboardingStatusEvent(a, cmd.Status, cmd.Comments, cmd.LoggedInUserId, cmd.CausedByContractId, updatedAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewUpdateOnboardingStatusEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) UpdateOrganizationOwner(ctx context.Context, cmd *command.UpdateOrganizationOwnerCommand) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.UpdateOrganizationOwner")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.GetTenant())
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.SetTag(tracing.SpanTagEntityId, cmd.ObjectID)
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "command", cmd)

	updatedAt := utils.Now()

	event, err := events.NewOrganizationOwnerUpdateEvent(a, cmd.OwnerUserId, cmd.ActorUserId, cmd.OrganizationId, updatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationOwnerUpdateEvent")
	}

	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: cmd.LoggedInUserId,
		App:    cmd.AppSource,
	})

	return a.Apply(event)
}
