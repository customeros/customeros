package aggregate

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organizationEvents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	OrganizationAggregateType eventstore.AggregateType = "organization"
)

type OrganizationAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Organization *model.Organization
}

func NewOrganizationAggregateWithTenantAndID(tenant, id string) *OrganizationAggregate {
	organizationAggregate := OrganizationAggregate{}
	organizationAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(OrganizationAggregateType, tenant, id)
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

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	socialId := request.SocialId
	if request.Url != "" && socialId == "" {
		if existingSocialId := a.Organization.GetSocialIdForUrl(request.Url); existingSocialId != "" {
			socialId = existingSocialId
		}
	}
	socialId = utils.NewUUIDIfEmpty(socialId)

	event, err := organizationEvents.NewOrganizationAddSocialEvent(a, socialId, request.Url, request.Alias, request.ExternalId, request.FollowersCount, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewOrganizationAddSocialEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	event, err := organizationEvents.NewOrganizationRemoveSocialEvent(a, socialId, request.Url)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationRemoveSocialEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	sourceFields := common.Source{}
	sourceFields.FromGrpc(request.SourceFields)
	sourceFields.SetDefaultValues()

	locationDtls := common.Location{
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

	event, err := organizationEvents.NewOrganizationAddLocationEvent(a, locationId, locationDtls, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewOrganizationAddLocationEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
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

	unlinkDomainEvent, err := organizationEvents.NewOrganizationUnlinkDomainEvent(a, request.Domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewOrganizationUnlinkDomainEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&unlinkDomainEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(unlinkDomainEvent)
}

func (a *OrganizationAggregate) When(event eventstore.Event) error {

	switch event.GetEventType() {
	case organizationEvents.OrganizationCreateV1:
		return a.onOrganizationCreate(event)
	case organizationEvents.OrganizationUpdateV1:
		return a.onOrganizationUpdate(event)
	case organizationEvents.OrganizationPhoneNumberLinkV1:
		return a.onPhoneNumberLink(event)
	case organizationEvents.OrganizationEmailLinkV1,
		organizationEvents.OrganizationEmailUnlinkV1:
		return nil
	case organizationEvents.OrganizationLocationLinkV1:
		return a.onLocationLink(event)
	case organizationEvents.OrganizationLinkDomainV1:
		return nil
	case organizationEvents.OrganizationUnlinkDomainV1:
		return nil
	case organizationEvents.OrganizationAddSocialV1:
		return a.onAddSocial(event)
	case organizationEvents.OrganizationRemoveSocialV1:
		return a.onRemoveSocial(event)
	case organizationEvents.OrganizationShowV1:
		return a.onShow(event)
	case organizationEvents.OrganizationUpsertCustomFieldV1:
		return a.onUpsertCustomField(event)
	case organizationEvents.OrganizationAddParentV1:
		return a.onAddParent(event)
	case organizationEvents.OrganizationRemoveParentV1:
		return a.onRemoveParent(event)
	case organizationEvents.OrganizationUpdateOnboardingStatusV1:
		return a.onOnboardingStatusUpdate(event)
	case organizationEvents.OrganizationUpdateOwnerV1:
		return a.onOrganizationOwnerUpdate(event)
	case organizationEvents.OrganizationCreateBillingProfileV1:
		return a.onCreateBillingProfile(event)
	case organizationEvents.OrganizationUpdateBillingProfileV1:
		return a.onUpdateBillingProfile(event)
	case organizationEvents.OrganizationEmailLinkToBillingProfileV1:
		return a.onEmailLinkToBillingProfile(event)
	case organizationEvents.OrganizationEmailUnlinkFromBillingProfileV1:
		return a.onEmailUnlinkFromBillingProfile(event)
	case organizationEvents.OrganizationLocationLinkToBillingProfileV1:
		return a.onLocationLinkToBillingProfile(event)
	case organizationEvents.OrganizationLocationUnlinkFromBillingProfileV1:
		return a.onLocationUnlinkFromBillingProfile(event)
	case organizationEvents.OrganizationAddLocationV1:
		return a.onAddLocation(event)
	case organizationEvents.OrganizationUpdateRenewalLikelihoodV1,
		organizationEvents.OrganizationUpdateRenewalForecastV1,
		organizationEvents.OrganizationUpdateBillingDetailsV1,
		organizationEvents.OrganizationRequestRenewalForecastV1,
		organizationEvents.OrganizationRequestNextCycleDateV1,
		organizationEvents.OrganizationRefreshLastTouchpointV1,
		organizationEvents.OrganizationRefreshArrV1,
		organizationEvents.OrganizationRefreshDerivedDataV1,
		organizationEvents.OrganizationRefreshRenewalSummaryV1,
		organizationEvents.OrganizationRequestScrapeByWebsiteV1,
		organizationEvents.OrganizationUpdateOwnerNotificationV1,
		organizationEvents.OrganizationRequestEnrichV1,
		organizationEvents.OrganizationHideV1,
		organizationEvents.OrganizationAddTagV1,
		organizationEvents.OrganizationRemoveTagV1:
		return nil
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
	var eventData organizationEvents.OrganizationCreateEvent
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
	a.Organization.IcpFit = eventData.IcpFit
	a.Organization.Source = common.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Organization.CreatedAt = eventData.CreatedAt
	a.Organization.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Organization.ExternalSystems = []common.ExternalSystem{eventData.ExternalSystem}
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
	var eventData organizationEvents.OrganizationUpdateEvent
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
		if eventData.UpdateIcpFit() {
			a.Organization.IcpFit = eventData.IcpFit
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
		if eventData.UpdateIcpFit() {
			a.Organization.IcpFit = eventData.IcpFit
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
	var eventData organizationEvents.OrganizationLinkPhoneNumberEvent
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

func (a *OrganizationAggregate) onLocationLink(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationLinkLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.LocationIds = utils.AddToListIfNotExists(a.Organization.LocationIds, eventData.LocationId)
	return nil
}

func (a *OrganizationAggregate) onAddSocial(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationAddSocialEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Socials == nil {
		a.Organization.Socials = make(map[string]common.Social)
	}
	a.Organization.Socials[eventData.SocialId] = common.Social{
		Url:            eventData.Url,
		Alias:          eventData.Alias,
		ExternalId:     eventData.ExternalId,
		FollowersCount: eventData.FollowersCount,
	}
	return nil
}

func (a *OrganizationAggregate) onRemoveSocial(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationAddSocialEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Socials == nil {
		a.Organization.Socials = make(map[string]common.Social)
	}
	delete(a.Organization.Socials, eventData.SocialId)
	return nil
}

func (a *OrganizationAggregate) onAddLocation(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationAddLocationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Organization.Locations == nil {
		a.Organization.Locations = make(map[string]common.Location)
	}
	a.Organization.Locations[eventData.LocationId] = common.Location{
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

func (a *OrganizationAggregate) onShow(event eventstore.Event) error {
	var eventData organizationEvents.ShowOrganizationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Organization.Hide = false
	return nil
}

func (a *OrganizationAggregate) onUpsertCustomField(event eventstore.Event) error {
	var eventData organizationEvents.OrganizationUpsertCustomField
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
			Source: common.Source{
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
	var eventData organizationEvents.OrganizationAddParentEvent
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
	var eventData organizationEvents.OrganizationRemoveParentEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	delete(a.Organization.ParentOrganizations, eventData.ParentOrganizationId)
	return nil
}

func (a *OrganizationAggregate) onOnboardingStatusUpdate(event eventstore.Event) error {
	var eventData organizationEvents.UpdateOnboardingStatusEvent
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
	var eventData organizationEvents.OrganizationOwnerUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	// do nothing
	return nil
}

func (a *OrganizationAggregate) onCreateBillingProfile(event eventstore.Event) error {
	var eventData organizationEvents.BillingProfileCreateEvent
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
	var eventData organizationEvents.BillingProfileUpdateEvent
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
	var eventData organizationEvents.LinkEmailToBillingProfileEvent
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
	var eventData organizationEvents.UnlinkEmailFromBillingProfileEvent
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
	var eventData organizationEvents.LinkLocationToBillingProfileEvent
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
	var eventData organizationEvents.UnlinkLocationFromBillingProfileEvent
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
