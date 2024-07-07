package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/aggregate"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	orgplanevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/events"
	orgplanmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	OrganizationAggregateType eventstore.AggregateType = "organization"
)

type OrganizationAggregate struct {
	*aggregate.CommonTenantIdAggregate
	Organization *model.Organization
}

func NewOrganizationAggregateWithTenantAndID(tenant, id string) *OrganizationAggregate {
	organizationAggregate := OrganizationAggregate{}
	organizationAggregate.CommonTenantIdAggregate = aggregate.NewCommonAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
	organizationAggregate.SetWhen(organizationAggregate.When)
	organizationAggregate.Organization = &model.Organization{}
	organizationAggregate.Tenant = tenant

	return &organizationAggregate
}

func (a *OrganizationAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *organizationpb.UnLinkDomainFromOrganizationGrpcRequest:
		return nil, a.unlinkDomain(ctx, r)
	case *organizationpb.OrganizationAddTagGrpcRequest:
		return nil, a.addTag(ctx, r)
	case *organizationpb.OrganizationRemoveTagGrpcRequest:
		return nil, a.removeTag(ctx, r)
	case *organizationpb.AddSocialGrpcRequest:
		return a.addSocial(ctx, r)
	case *organizationpb.RemoveSocialGrpcRequest:
		return nil, a.removeSocial(ctx, r)
	case *organizationpb.OrganizationAddLocationGrpcRequest:
		return a.addLocation(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *OrganizationAggregate) addSocial(ctx context.Context, request *organizationpb.AddSocialGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addSocial")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	sourceFields := cmnmod.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	socialId := request.SocialId
	if request.Url != "" && socialId == "" {
		if existingSocialId := a.Organization.GetSocialIdForUrl(request.Url); existingSocialId != "" {
			socialId = existingSocialId
		}
	}
	socialId = utils.NewUUIDIfEmpty(socialId)

	event, err := events.NewOrganizationAddSocialEvent(a, socialId, request.Url, request.Alias, request.ExternalId, request.FollowersCount, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return socialId, a.Apply(event)
}

func (a *OrganizationAggregate) removeSocial(ctx context.Context, request *organizationpb.RemoveSocialGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.removeSocial")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	socialId := request.SocialId
	if socialId == "" {
		if existingSocialId := a.Organization.GetSocialIdForUrl(request.Url); existingSocialId != "" {
			socialId = existingSocialId
		}
	}

	event, err := events.NewOrganizationRemoveSocialEvent(a, socialId, request.Url)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRemoveSocialEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(event)
}

func (a *OrganizationAggregate) addLocation(ctx context.Context, request *organizationpb.OrganizationAddLocationGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	createdAtNotNil := utils.IfNotNilTimeWithDefault(utils.TimestampProtoToTimePtr(request.CreatedAt), utils.Now())

	sourceFields := cmnmod.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	locationDtls := cmnmod.Location{
		Name:          request.LocationDetails.Name,
		RawAddress:    request.LocationDetails.RawAddress,
		Country:       request.LocationDetails.Country,
		CountryCodeA2: request.LocationDetails.CountryCodeA2,
		CountryCodeA3: request.LocationDetails.CountryCodeA3,
		Region:        request.LocationDetails.Region,
		Locality:      request.LocationDetails.Locality,
		AddressLine1:  request.LocationDetails.AddressLine1,
		AddressLine2:  request.LocationDetails.AddressLine2,
		Street:        request.LocationDetails.Street,
		HouseNumber:   request.LocationDetails.HouseNumber,
		ZipCode:       request.LocationDetails.ZipCode,
		PostalCode:    request.LocationDetails.PostalCode,
		AddressType:   request.LocationDetails.AddressType,
		Commercial:    request.LocationDetails.Commercial,
		Predirection:  request.LocationDetails.Predirection,
		PlusFour:      request.LocationDetails.PlusFour,
		TimeZone:      request.LocationDetails.TimeZone,
		UtcOffset:     request.LocationDetails.UtcOffset,
		Latitude:      utils.ParseStringToFloat(request.LocationDetails.Latitude),
		Longitude:     utils.ParseStringToFloat(request.LocationDetails.Longitude),
	}

	locationId := request.LocationId
	if locationId == "" && !locationDtls.IsEmpty() {
		if existingLocaitonId := a.Organization.GetLocationIdForDetails(locationDtls); existingLocaitonId != "" {
			locationId = existingLocaitonId
		}
	}
	locationId = utils.NewUUIDIfEmpty(locationId)

	event, err := events.NewOrganizationAddLocationEvent(a, locationId, locationDtls, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewOrganizationAddLocationEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&event, span, aggregate.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return locationId, a.Apply(event)
}

func (a *OrganizationAggregate) unlinkDomain(ctx context.Context, request *organizationpb.UnLinkDomainFromOrganizationGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.unlinkDomain")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	unlinkDomainEvent, err := events.NewOrganizationUnlinkDomainEvent(a, request.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUnlinkDomainEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&unlinkDomainEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(unlinkDomainEvent)
}

func (a *OrganizationAggregate) addTag(ctx context.Context, request *organizationpb.OrganizationAddTagGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.addTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	addTagEvent, err := events.NewOrganizationAddTagEvent(a, request.TagId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationAddTagEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&addTagEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(addTagEvent)
}

func (a *OrganizationAggregate) removeTag(ctx context.Context, request *organizationpb.OrganizationRemoveTagGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OrganizationAggregate.removeTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	removeTagEvent, err := events.NewOrganizationRemoveTagEvent(a, request.TagId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRemoveTagEvent")
	}
	aggregate.EnrichEventWithMetadataExtended(&removeTagEvent, span, aggregate.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(removeTagEvent)
}

func (a *OrganizationAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {
	case events.OrganizationCreateV1:
		return a.onOrganizationCreate(event)
	case events.OrganizationUpdateV1:
		return a.onOrganizationUpdate(event)
	case events.OrganizationPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case events.OrganizationEmailLinkV1:
		return a.onEmailLink(event)
	case events.OrganizationLocationLinkV1:
		return a.onLocationLink(event)
	case events.OrganizationLinkDomainV1:
		return a.onDomainLink(event)
	case events.OrganizationUnlinkDomainV1:
		return a.onDomainUnlink(event)
	case events.OrganizationAddSocialV1:
		return a.onAddSocial(event)
	case events.OrganizationRemoveSocialV1:
		return a.onRemoveSocial(event)
	case events.OrganizationHideV1:
		return a.onHide(event)
	case events.OrganizationShowV1:
		return a.onShow(event)
	case events.OrganizationUpsertCustomFieldV1:
		return a.onUpsertCustomField(event)
	case events.OrganizationAddParentV1:
		return a.onAddParent(event)
	case events.OrganizationRemoveParentV1:
		return a.onRemoveParent(event)
	case events.OrganizationUpdateOnboardingStatusV1:
		return a.onOnboardingStatusUpdate(event)
	case events.OrganizationUpdateOwnerV1:
		return a.onOrganizationOwnerUpdate(event)
	case events.OrganizationCreateBillingProfileV1:
		return a.onCreateBillingProfile(event)
	case events.OrganizationUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(event)
	case events.OrganizationEmailLinkToBillingProfileV1:
		return a.onEmailLinkToBillingProfile(event)
	case events.OrganizationEmailUnlinkFromBillingProfileV1:
		return a.onEmailUnlinkFromBillingProfile(event)
	case events.OrganizationLocationLinkToBillingProfileV1:
		return a.onLocationLinkToBillingProfile(event)
	case events.OrganizationLocationUnlinkFromBillingProfileV1:
		return a.onLocationUnlinkFromBillingProfile(event)
	case events.OrganizationAddTagV1:
		return a.onOrganizationAddTag(event)
	case events.OrganizationRemoveTagV1:
		return a.onOrganizationRemoveTag(event)
	case events.OrganizationUpdateRenewalLikelihoodV1,
		events.OrganizationUpdateRenewalForecastV1,
		events.OrganizationUpdateBillingDetailsV1,
		events.OrganizationRequestRenewalForecastV1,
		events.OrganizationRequestNextCycleDateV1,
		events.OrganizationRefreshLastTouchpointV1,
		events.OrganizationRefreshArrV1,
		events.OrganizationRefreshDerivedDataV1,
		events.OrganizationRefreshRenewalSummaryV1,
		events.OrganizationRequestScrapeByWebsiteV1,
		events.OrganizationUpdateOwnerNotificationV1,
		events.OrganizationRequestEnrichV1:
		return nil
	case orgplanevents.OrganizationPlanCreateV1:
		return a.onOrganizationPlanCreate(event)
	case orgplanevents.OrganizationPlanUpdateV1:
		return a.onOrganizationPlanUpdate(event)
	case orgplanevents.OrganizationPlanMilestoneCreateV1:
		return a.onOrganizationPlanMilestoneCreate(event)
	case orgplanevents.OrganizationPlanMilestoneUpdateV1:
		return a.onOrganizationPlanMilestoneUpdate(event)
	case orgplanevents.OrganizationPlanMilestoneReorderV1:
		return a.onOrganizationPlanMilestoneReorder(event)
	case events.OrganizationAddLocationV1:
		return a.onAddLocation(event)
	default:
		if strings.HasPrefix(event.GetEventType(), constants.EsInternalStreamPrefix) {
			return nil
		}
		span, _ := opentracing.StartSpanFromContext(context.Background(), "OrganizationAggregate.When")
		defer span.Finish()
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		tracing.TraceErr(span, eventstore.ErrInvalidEventType)
		return err
	}
}

func (a *OrganizationAggregate) onOrganizationCreate(event eventstore.Event) error {
	var eventData events.OrganizationCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Name = eventData.Name
	a.Organization.Hide = eventData.Hide
	a.Organization.Description = eventData.Description
	a.Organization.Website = eventData.Website
	a.Organization.Industry = eventData.Industry
	a.Organization.SubIndustry = eventData.SubIndustry
	a.Organization.IndustryGroup = eventData.IndustryGroup
	a.Organization.TargetAudience = eventData.TargetAudience
	a.Organization.ValueProposition = eventData.ValueProposition
	a.Organization.LastFundingRound = eventData.LastFundingRound
	a.Organization.LastFundingAmount = eventData.LastFundingAmount
	a.Organization.ReferenceId = eventData.ReferenceId
	a.Organization.Note = eventData.Note
	a.Organization.IsPublic = eventData.IsPublic
	a.Organization.Relationship = eventData.Relationship
	a.Organization.Stage = eventData.Stage
	a.Organization.Employees = eventData.Employees
	a.Organization.Market = eventData.Market
	a.Organization.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Organization.CreatedAt = eventData.CreatedAt
	a.Organization.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Organization.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	a.Organization.YearFounded = eventData.YearFounded
	a.Organization.Headquarters = eventData.Headquarters
	a.Organization.EmployeeGrowthRate = eventData.EmployeeGrowthRate
	a.Organization.SlackChannelId = eventData.SlackChannelId
	a.Organization.LogoUrl = eventData.LogoUrl
	a.Organization.IconUrl = eventData.IconUrl
	a.Organization.LeadSource = eventData.LeadSource
	return nil
}

func (a *OrganizationAggregate) onOrganizationUpdate(event eventstore.Event) error {
	var eventData events.OrganizationUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// Update only if the source of truth is 'openline' or the new source matches the source of truth
	if eventData.Source == constants.SourceOpenline {
		a.Organization.Source.SourceOfTruth = eventData.Source
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt

	if !eventData.Hide {
		a.Organization.Hide = false
	}
	if eventData.Source != a.Organization.Source.SourceOfTruth && a.Organization.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Organization.Name == "" && eventData.UpdateName() {
			a.Organization.Name = eventData.Name
		}
		if a.Organization.Description == "" && eventData.UpdateDescription() {
			a.Organization.Description = eventData.Description
		}
		if a.Organization.Website == "" && eventData.UpdateWebsite() {
			a.Organization.Website = eventData.Website
		}
		if a.Organization.Industry == "" && eventData.UpdateIndustry() {
			a.Organization.Industry = eventData.Industry
		}
		if a.Organization.SubIndustry == "" && eventData.UpdateSubIndustry() {
			a.Organization.SubIndustry = eventData.SubIndustry
		}
		if a.Organization.IndustryGroup == "" && eventData.UpdateIndustryGroup() {
			a.Organization.IndustryGroup = eventData.IndustryGroup
		}
		if a.Organization.TargetAudience == "" && eventData.UpdateTargetAudience() {
			a.Organization.TargetAudience = eventData.TargetAudience
		}
		if a.Organization.ValueProposition == "" && eventData.UpdateValueProposition() {
			a.Organization.ValueProposition = eventData.ValueProposition
		}
		if a.Organization.LastFundingRound == "" && eventData.UpdateLastFundingRound() {
			a.Organization.LastFundingRound = eventData.LastFundingRound
		}
		if a.Organization.LastFundingAmount == "" && eventData.UpdateLastFundingAmount() {
			a.Organization.LastFundingAmount = eventData.LastFundingAmount
		}
		if a.Organization.ReferenceId == "" && eventData.UpdateReferenceId() {
			a.Organization.ReferenceId = eventData.ReferenceId
		}
		if a.Organization.Note == "" && eventData.UpdateNote() {
			a.Organization.Note = eventData.Note
		}
		if a.Organization.Employees == 0 && eventData.UpdateEmployees() {
			a.Organization.Employees = eventData.Employees
		}
		if a.Organization.Market == "" && eventData.UpdateMarket() {
			a.Organization.Market = eventData.Market
		}
		if a.Organization.Relationship == "" && eventData.UpdateRelationship() {
			a.Organization.Relationship = eventData.Relationship
		}
		if a.Organization.Stage == "" && eventData.UpdateStage() {
			a.Organization.Stage = eventData.Stage
		}
		if a.Organization.YearFounded == nil && eventData.UpdateYearFounded() {
			a.Organization.YearFounded = eventData.YearFounded
		}
		if a.Organization.Headquarters == "" && eventData.UpdateHeadquarters() {
			a.Organization.Headquarters = eventData.Headquarters
		}
		if a.Organization.EmployeeGrowthRate == "" && eventData.UpdateEmployeeGrowthRate() {
			a.Organization.EmployeeGrowthRate = eventData.EmployeeGrowthRate
		}
		if a.Organization.SlackChannelId == "" && eventData.UpdateSlackChannelId() {
			a.Organization.SlackChannelId = eventData.SlackChannelId
		}
		if a.Organization.LogoUrl == "" && eventData.UpdateLogoUrl() {
			a.Organization.LogoUrl = eventData.LogoUrl
		}
		if a.Organization.IconUrl == "" && eventData.UpdateIconUrl() {
			a.Organization.IconUrl = eventData.IconUrl
		}
	} else {
		if eventData.UpdateIsPublic() {
			a.Organization.IsPublic = eventData.IsPublic
		}
		if eventData.UpdateRelationship() {
			a.Organization.Relationship = eventData.Relationship
		}
		if eventData.UpdateStage() {
			a.Organization.Stage = eventData.Stage
		}
		if eventData.UpdateHide() {
			a.Organization.Hide = eventData.Hide
		}
		if eventData.UpdateName() {
			a.Organization.Name = eventData.Name
		}
		if eventData.UpdateDescription() {
			a.Organization.Description = eventData.Description
		}
		if eventData.UpdateWebsite() {
			a.Organization.Website = eventData.Website
		}
		if eventData.UpdateIndustry() {
			a.Organization.Industry = eventData.Industry
		}
		if eventData.UpdateSubIndustry() {
			a.Organization.SubIndustry = eventData.SubIndustry
		}
		if eventData.UpdateIndustryGroup() {
			a.Organization.IndustryGroup = eventData.IndustryGroup
		}
		if eventData.UpdateTargetAudience() {
			a.Organization.TargetAudience = eventData.TargetAudience
		}
		if eventData.UpdateValueProposition() {
			a.Organization.ValueProposition = eventData.ValueProposition
		}
		if eventData.UpdateLastFundingRound() {
			a.Organization.LastFundingRound = eventData.LastFundingRound
		}
		if eventData.UpdateLastFundingAmount() {
			a.Organization.LastFundingAmount = eventData.LastFundingAmount
		}
		if eventData.UpdateReferenceId() {
			a.Organization.ReferenceId = eventData.ReferenceId
		}
		if eventData.UpdateNote() {
			a.Organization.Note = eventData.Note
		}
		if eventData.UpdateEmployees() {
			a.Organization.Employees = eventData.Employees
		}
		if eventData.UpdateMarket() {
			a.Organization.Market = eventData.Market
		}
		if eventData.UpdateYearFounded() {
			a.Organization.YearFounded = eventData.YearFounded
		}
		if eventData.UpdateHeadquarters() {
			a.Organization.Headquarters = eventData.Headquarters
		}
		if eventData.UpdateEmployeeGrowthRate() {
			a.Organization.EmployeeGrowthRate = eventData.EmployeeGrowthRate
		}
		if eventData.UpdateSlackChannelId() {
			a.Organization.SlackChannelId = eventData.SlackChannelId
		}
		if eventData.UpdateLogoUrl() {
			a.Organization.LogoUrl = eventData.LogoUrl
		}
		if eventData.UpdateIconUrl() {
			a.Organization.IconUrl = eventData.IconUrl
		}
	}
	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Organization.ExternalSystems {
			if externalSystem.ExternalSystemId == eventData.ExternalSystem.ExternalSystemId &&
				externalSystem.ExternalId == eventData.ExternalSystem.ExternalId {
				found = true
				externalSystem.ExternalUrl = eventData.ExternalSystem.ExternalUrl
				externalSystem.SyncDate = eventData.ExternalSystem.SyncDate
				externalSystem.ExternalSource = eventData.ExternalSystem.ExternalSource
				if eventData.ExternalSystem.ExternalIdSecond != "" {
					externalSystem.ExternalIdSecond = eventData.ExternalSystem.ExternalIdSecond
				}
			}
		}
		if !found {
			a.Organization.ExternalSystems = append(a.Organization.ExternalSystems, eventData.ExternalSystem)
		}
	}
	return nil
}

func (a *OrganizationAggregate) onPhoneNumberLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.PhoneNumbers == nil {
		a.Organization.PhoneNumbers = make(map[string]model.OrganizationPhoneNumber)
	}
	a.Organization.PhoneNumbers[eventData.PhoneNumberId] = model.OrganizationPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onEmailLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkEmailEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Emails == nil {
		a.Organization.Emails = make(map[string]model.OrganizationEmail)
	}
	a.Organization.Emails[eventData.EmailId] = model.OrganizationEmail{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Organization.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *OrganizationAggregate) onDomainLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkDomainEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Domains == nil {
		a.Organization.Domains = []string{}
	}
	if !utils.Contains(a.Organization.Domains, strings.TrimSpace(eventData.Domain)) {
		a.Organization.Domains = append(a.Organization.Domains, strings.TrimSpace(eventData.Domain))
	}
	return nil
}

func (a *OrganizationAggregate) onDomainUnlink(event eventstore.Event) error {
	var eventData events.OrganizationLinkDomainEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Domains == nil {
		a.Organization.Domains = []string{}
	}
	utils.RemoveFromList(a.Organization.Domains, eventData.Domain)
	return nil
}

func (a *OrganizationAggregate) onLocationLink(event eventstore.Event) error {
	var eventData events.OrganizationLinkLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.LocationIds = utils.AddToListIfNotExists(a.Organization.LocationIds, eventData.LocationId)
	return nil
}

func (a *OrganizationAggregate) onAddSocial(event eventstore.Event) error {
	var eventData events.OrganizationAddSocialEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Socials == nil {
		a.Organization.Socials = make(map[string]cmnmod.Social)
	}
	a.Organization.Socials[eventData.SocialId] = cmnmod.Social{
		Url:            eventData.Url,
		Alias:          eventData.Alias,
		ExternalId:     eventData.ExternalId,
		FollowersCount: eventData.FollowersCount,
	}
	return nil
}

func (a *OrganizationAggregate) onRemoveSocial(event eventstore.Event) error {
	var eventData events.OrganizationAddSocialEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Socials == nil {
		a.Organization.Socials = make(map[string]cmnmod.Social)
	}
	delete(a.Organization.Socials, eventData.SocialId)
	return nil
}

func (a *OrganizationAggregate) onAddLocation(event eventstore.Event) error {
	var eventData events.OrganizationAddLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Locations == nil {
		a.Organization.Locations = make(map[string]cmnmod.Location)
	}
	a.Organization.Locations[eventData.LocationId] = cmnmod.Location{
		Name:          eventData.Name,
		RawAddress:    eventData.RawAddress,
		Country:       eventData.Country,
		CountryCodeA2: eventData.CountryCodeA2,
		CountryCodeA3: eventData.CountryCodeA3,
		Region:        eventData.Region,
		Locality:      eventData.Locality,
		AddressLine1:  eventData.AddressLine1,
		AddressLine2:  eventData.AddressLine2,
		Street:        eventData.Street,
		HouseNumber:   eventData.HouseNumber,
		ZipCode:       eventData.ZipCode,
		PostalCode:    eventData.PostalCode,
		AddressType:   eventData.AddressType,
		Commercial:    eventData.Commercial,
		Predirection:  eventData.Predirection,
		PlusFour:      eventData.PlusFour,
		TimeZone:      eventData.TimeZone,
		UtcOffset:     eventData.UtcOffset,
		Latitude:      eventData.Latitude,
		Longitude:     eventData.Longitude,
	}
	return nil
}

func (a *OrganizationAggregate) onHide(event eventstore.Event) error {
	var eventData events.HideOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Hide = true
	return nil
}

func (a *OrganizationAggregate) onShow(event eventstore.Event) error {
	var eventData events.ShowOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Hide = false
	return nil
}

func (a *OrganizationAggregate) onUpsertCustomField(event eventstore.Event) error {
	var eventData events.OrganizationUpsertCustomField
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.CustomFields == nil {
		a.Organization.CustomFields = make(map[string]model.CustomField)
	}

	if val, ok := a.Organization.CustomFields[eventData.CustomFieldId]; ok {
		val.Source.SourceOfTruth = eventData.SourceOfTruth
		val.UpdatedAt = eventData.UpdatedAt
		val.CustomFieldValue = eventData.CustomFieldValue
		val.Name = eventData.CustomFieldName
	} else {
		a.Organization.CustomFields[eventData.CustomFieldId] = model.CustomField{
			Source: cmnmod.Source{
				Source:        eventData.Source,
				SourceOfTruth: eventData.SourceOfTruth,
				AppSource:     eventData.AppSource,
			},
			CreatedAt:           eventData.CreatedAt,
			UpdatedAt:           eventData.UpdatedAt,
			Id:                  eventData.CustomFieldId,
			TemplateId:          eventData.TemplateId,
			Name:                eventData.CustomFieldName,
			CustomFieldDataType: model.CustomFieldDataType(eventData.CustomFieldDataType),
			CustomFieldValue:    eventData.CustomFieldValue,
		}
	}
	return nil
}

func (a *OrganizationAggregate) onAddParent(event eventstore.Event) error {
	var eventData events.OrganizationAddParentEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.ParentOrganizations == nil {
		a.Organization.ParentOrganizations = make(map[string]model.ParentOrganization)
	}
	a.Organization.ParentOrganizations[eventData.ParentOrganizationId] = model.ParentOrganization{
		OrganizationId: eventData.ParentOrganizationId,
		Type:           eventData.Type,
	}
	return nil
}

func (a *OrganizationAggregate) onRemoveParent(event eventstore.Event) error {
	var eventData events.OrganizationRemoveParentEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	delete(a.Organization.ParentOrganizations, eventData.ParentOrganizationId)
	return nil
}

func (a *OrganizationAggregate) onOnboardingStatusUpdate(event eventstore.Event) error {
	var eventData events.UpdateOnboardingStatusEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Organization.OnboardingDetails = model.OnboardingDetails{
		Status:    eventData.Status,
		Comments:  eventData.Comments,
		UpdatedAt: eventData.UpdatedAt,
	}

	return nil
}

func (a *OrganizationAggregate) onOrganizationOwnerUpdate(event eventstore.Event) error {
	var eventData events.OrganizationOwnerUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// do nothing
	return nil
}

func (a *OrganizationAggregate) onCreateBillingProfile(event eventstore.Event) error {
	var eventData events.BillingProfileCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}

	a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{
		Id:           eventData.BillingProfileId,
		LegalName:    eventData.LegalName,
		TaxId:        eventData.TaxId,
		CreatedAt:    eventData.CreatedAt,
		UpdatedAt:    eventData.UpdatedAt,
		SourceFields: eventData.SourceFields,
	}

	return nil
}

func (a *OrganizationAggregate) onUpdateBillingProfile(event eventstore.Event) error {
	var eventData events.BillingProfileUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	if eventData.UpdateLegalName() {
		billingProfile.LegalName = eventData.LegalName
	}
	if eventData.UpdateTaxId() {
		billingProfile.TaxId = eventData.TaxId
	}
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onEmailLinkToBillingProfile(event eventstore.Event) error {
	var eventData events.LinkEmailToBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	if eventData.Primary {
		billingProfile.PrimaryEmailId = eventData.EmailId
		billingProfile.EmailIds = utils.RemoveFromList(billingProfile.EmailIds, eventData.EmailId)
	} else {
		billingProfile.EmailIds = utils.AddToListIfNotExists(billingProfile.EmailIds, eventData.EmailId)
		if billingProfile.PrimaryEmailId == eventData.EmailId {
			billingProfile.PrimaryEmailId = ""
		}
	}
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onEmailUnlinkFromBillingProfile(event eventstore.Event) error {
	var eventData events.UnlinkEmailFromBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}
	if billingProfile.PrimaryEmailId == eventData.EmailId {
		billingProfile.PrimaryEmailId = ""
	}
	billingProfile.EmailIds = utils.RemoveFromList(billingProfile.EmailIds, eventData.EmailId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onLocationLinkToBillingProfile(event eventstore.Event) error {
	var eventData events.LinkLocationToBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}

	billingProfile.LocationIds = utils.AddToListIfNotExists(billingProfile.LocationIds, eventData.LocationId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onLocationUnlinkFromBillingProfile(event eventstore.Event) error {
	var eventData events.UnlinkLocationFromBillingProfileEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.BillingProfiles == nil {
		a.Organization.BillingProfiles = make(map[string]model.BillingProfile)
	}
	billingProfile, ok := a.Organization.BillingProfiles[eventData.BillingProfileId]
	if !ok {
		a.Organization.BillingProfiles[eventData.BillingProfileId] = model.BillingProfile{}
		billingProfile = a.Organization.BillingProfiles[eventData.BillingProfileId]
	}
	billingProfile.LocationIds = utils.RemoveFromList(billingProfile.LocationIds, eventData.LocationId)
	billingProfile.UpdatedAt = eventData.UpdatedAt
	a.Organization.BillingProfiles[eventData.BillingProfileId] = billingProfile

	return nil
}

func (a *OrganizationAggregate) onOrganizationPlanCreate(event eventstore.Event) error {
	var eventData orgplanevents.OrganizationPlanCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.OrganizationPlans == nil {
		a.Organization.OrganizationPlans = map[string]orgplanmodel.OrganizationPlan{}
	}

	a.Organization.OrganizationPlans[eventData.OrganizationPlanId] = orgplanmodel.OrganizationPlan{
		ID:           eventData.OrganizationPlanId,
		Name:         eventData.Name,
		SourceFields: eventData.SourceFields,
		CreatedAt:    eventData.CreatedAt,
		UpdatedAt:    eventData.CreatedAt,
		MasterPlanId: eventData.MasterPlanId,
	}

	return nil
}

func (a *OrganizationAggregate) onOrganizationPlanUpdate(event eventstore.Event) error {
	var eventData orgplanevents.OrganizationPlanUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.OrganizationPlans == nil {
		a.Organization.OrganizationPlans = map[string]orgplanmodel.OrganizationPlan{}
	}

	op := orgplanmodel.OrganizationPlan{
		ID:        eventData.OrganizationPlanId,
		Name:      eventData.Name,
		UpdatedAt: eventData.UpdatedAt,
	}

	a.Organization.OrganizationPlans[eventData.OrganizationPlanId] = op

	return nil
}

func (a *OrganizationAggregate) onOrganizationPlanMilestoneCreate(event eventstore.Event) error {
	var eventData orgplanevents.OrganizationPlanMilestoneCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.OrganizationPlans == nil {
		a.Organization.OrganizationPlans = map[string]orgplanmodel.OrganizationPlan{}
	}

	mstone := orgplanmodel.OrganizationPlanMilestone{
		ID:        eventData.MilestoneId,
		Name:      eventData.Name,
		Order:     eventData.Order,
		CreatedAt: eventData.CreatedAt,
		DueDate:   eventData.DueDate,
		UpdatedAt: eventData.CreatedAt,
		Optional:  eventData.Optional,
		Items:     convertMilestoneItemsToObject(eventData.Items),
	}

	if a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones == nil {
		// First we get a "copy" of the entry; doing this because we can't modify the map entry directly
		if entry, ok := a.Organization.OrganizationPlans[eventData.OrganizationPlanId]; ok {

			// Then we modify the copy
			entry.Milestones = make(map[string]orgplanmodel.OrganizationPlanMilestone)

			// Then we reassign map entry
			a.Organization.OrganizationPlans[eventData.OrganizationPlanId] = entry
		}
	}
	a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[mstone.ID] = mstone

	return nil
}

func convertMilestoneItemsToObject(items []string) []orgplanmodel.OrganizationPlanMilestoneItem {
	milestoneItems := make([]orgplanmodel.OrganizationPlanMilestoneItem, len(items))
	for i, item := range items {
		milestoneItems[i] = orgplanmodel.OrganizationPlanMilestoneItem{
			Text:      item,
			UpdatedAt: time.Now(),
			Status:    orgplanmodel.TaskNotDone.String(),
		}
	}
	return milestoneItems
}

func (a *OrganizationAggregate) onOrganizationPlanMilestoneUpdate(evt eventstore.Event) error {
	var eventData orgplanevents.OrganizationPlanMilestoneUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.OrganizationPlans == nil {
		a.Organization.OrganizationPlans = map[string]orgplanmodel.OrganizationPlan{}
	}

	if a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones == nil {
		// First we get a "copy" of the entry; doing this because we can't modify the map entry directly
		if entry, ok := a.Organization.OrganizationPlans[eventData.OrganizationPlanId]; ok {

			// Then we modify the copy
			entry.Milestones = make(map[string]orgplanmodel.OrganizationPlanMilestone)

			// Then we reassign map entry
			a.Organization.OrganizationPlans[eventData.OrganizationPlanId] = entry
		}
	}
	if _, ok := a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[eventData.MilestoneId]; !ok {
		a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[eventData.MilestoneId] = orgplanmodel.OrganizationPlanMilestone{
			ID: eventData.MilestoneId,
		}
	}
	milestone := a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[eventData.MilestoneId]
	if eventData.UpdateName() {
		milestone.Name = eventData.Name
	}
	if eventData.UpdateOrder() {
		milestone.Order = eventData.Order
	}
	if eventData.UpdateDueDate() {
		milestone.DueDate = eventData.DueDate
	}
	if eventData.UpdateItems() {
		milestone.Items = eventData.Items
	}
	if eventData.UpdateOptional() {
		milestone.Optional = eventData.Optional
	}
	if eventData.UpdateRetired() {
		milestone.Retired = eventData.Retired
	}
	if eventData.UpdateStatusDetails() {
		milestone.StatusDetails = eventData.StatusDetails
	}

	a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[milestone.ID] = milestone

	return nil
}

func (a *OrganizationAggregate) onOrganizationPlanMilestoneReorder(evt eventstore.Event) error {
	var eventData orgplanevents.OrganizationPlanMilestoneReorderEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if a.Organization.OrganizationPlans == nil {
		a.Organization.OrganizationPlans = map[string]orgplanmodel.OrganizationPlan{}
	}

	if a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones == nil {
		// First we get a "copy" of the entry; doing this because we can't modify the map entry directly
		if entry, ok := a.Organization.OrganizationPlans[eventData.OrganizationPlanId]; ok {

			// Then we modify the copy
			entry.Milestones = make(map[string]orgplanmodel.OrganizationPlanMilestone)

			// Then we reassign map entry
			a.Organization.OrganizationPlans[eventData.OrganizationPlanId] = entry
		}
	}
	for i, milestoneId := range eventData.MilestoneIds {
		if milestone, ok := a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[milestoneId]; ok {
			milestone.Order = int64(i)
			a.Organization.OrganizationPlans[eventData.OrganizationPlanId].Milestones[milestoneId] = milestone
		}
	}

	return nil
}

func (a *OrganizationAggregate) onOrganizationAddTag(evt eventstore.Event) error {
	var eventData events.OrganizationAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Organization.TagIds = append(a.Organization.TagIds, eventData.TagId)
	a.Organization.TagIds = utils.RemoveDuplicates(a.Organization.TagIds)

	return nil
}

func (a *OrganizationAggregate) onOrganizationRemoveTag(evt eventstore.Event) error {
	var eventData events.OrganizationRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Organization.TagIds = utils.RemoveFromList(a.Organization.TagIds, eventData.TagId)

	return nil
}
