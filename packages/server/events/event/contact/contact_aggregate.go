package contact

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/openline-ai/openline-customer-os/packages/server/events/constants"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/event/contact/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	ContactAggregateType       eventstore.AggregateType = "contact"
	PARAM_REQUEST              string                   = "request"
	PARAM_REQUEST_SHOW_CONTACT string                   = "showContact"
	PARAM_REQUEST_HIDE_CONTACT string                   = "hideContact"
)

type ContactAggregate struct {
	*eventstore.CommonTenantIdAggregate
	Contact *Contact
}

func NewContactAggregateWithTenantAndID(tenant, id string) *ContactAggregate {
	contactAggregate := ContactAggregate{}
	contactAggregate.CommonTenantIdAggregate = eventstore.NewCommonAggregateWithTenantAndId(ContactAggregateType, tenant, id)
	contactAggregate.SetWhen(contactAggregate.When)
	contactAggregate.Contact = &Contact{}
	contactAggregate.Tenant = tenant
	return &contactAggregate
}

func (a *ContactAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactAggregate.HandleGRPCRequest")
	defer span.Finish()

	switch r := request.(type) {
	case *contactpb.ContactAddTagGrpcRequest:
		return nil, a.addTag(ctx, r)
	case *contactpb.ContactRemoveTagGrpcRequest:
		return nil, a.removeTag(ctx, r)
	case *contactpb.ContactAddSocialGrpcRequest:
		return a.addSocial(ctx, r)
	case *contactpb.ContactRemoveSocialGrpcRequest:
		return nil, a.removeSocial(ctx, r)
	case *contactpb.ContactAddLocationGrpcRequest:
		return a.addLocation(ctx, r)
	case *contactpb.UnLinkEmailFromContactGrpcRequest:
		return nil, a.unlinkEmail(ctx, r)
	case *contactpb.ContactIdGrpcRequest:
		requestType := ""
		if params != nil {
			if _, ok := params[PARAM_REQUEST]; ok {
				requestType = params[PARAM_REQUEST].(string)
			}
		}
		switch requestType {
		case PARAM_REQUEST_SHOW_CONTACT:
			return nil, a.showContact(ctx, r)
		case PARAM_REQUEST_HIDE_CONTACT:
			return nil, a.hideContact(ctx, r)
		}
		return nil, nil
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
}

func (a *ContactAggregate) addTag(ctx context.Context, request *contactpb.ContactAddTagGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.addTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	addTagEvent, err := event.NewContactAddTagEvent(a, request.TagId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactAddTagEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&addTagEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(addTagEvent)
}

func (a *ContactAggregate) removeTag(ctx context.Context, request *contactpb.ContactRemoveTagGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.removeTag")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	removeTagEvent, err := event.NewContactRemoveTagEvent(a, request.TagId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactRemoveTagEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&removeTagEvent, span, eventstore.EventMetadata{
		Tenant: a.Tenant,
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(removeTagEvent)
}

func (a *ContactAggregate) addSocial(ctx context.Context, request *contactpb.ContactAddSocialGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.addSocial")
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
		if existingSocialId := a.Contact.GetSocialIdForUrl(request.Url); existingSocialId != "" {
			socialId = existingSocialId
		}
	}
	socialId = utils.NewUUIDIfEmpty(socialId)

	addSocialEvent, err := event.NewContactAddSocialEvent(a, socialId, request.Url, request.Alias, request.ExternalId, request.FollowersCount, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewContactAddSocialEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&addSocialEvent, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return socialId, a.Apply(addSocialEvent)
}

func (a *ContactAggregate) removeSocial(ctx context.Context, request *contactpb.ContactRemoveSocialGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.removeSocial")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	socialId := request.SocialId
	if socialId == "" {
		if existingSocialId := a.Contact.GetSocialIdForUrl(request.Url); existingSocialId != "" {
			socialId = existingSocialId
		}
	}

	removeSocialEvent, err := event.NewContactRemoveSocialEvent(a, socialId, request.Url)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactRemoveSocialEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&removeSocialEvent, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(removeSocialEvent)
}

func (a *ContactAggregate) addLocation(ctx context.Context, request *contactpb.ContactAddLocationGrpcRequest) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.addLocation")
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
		if existingLocaitonId := a.Contact.GetLocationIdForDetails(locationDtls); existingLocaitonId != "" {
			locationId = existingLocaitonId
		}
	}
	locationId = utils.NewUUIDIfEmpty(locationId)

	event, err := event.NewContactAddLocationEvent(a, locationId, locationDtls, sourceFields, createdAtNotNil)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", errors.Wrap(err, "NewContactAddLocationEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&event, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    sourceFields.AppSource,
	})

	return locationId, a.Apply(event)
}

func (a *ContactAggregate) showContact(ctx context.Context, request *contactpb.ContactIdGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.showContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	showContactEvent, err := event.NewContactShowEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactShowEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&showContactEvent, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(showContactEvent)
}

func (a *ContactAggregate) hideContact(ctx context.Context, request *contactpb.ContactIdGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.showContact")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	hideContactEvent, err := event.NewContactHideEvent(a)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactHideEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&hideContactEvent, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(hideContactEvent)
}

func (a *ContactAggregate) unlinkEmail(ctx context.Context, request *contactpb.UnLinkEmailFromContactGrpcRequest) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ContactAggregate.unlinkEmail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, a.Tenant)
	span.SetTag(tracing.SpanTagAggregateId, a.GetID())
	span.LogFields(log.Int64("aggregateVersion", a.GetVersion()))
	tracing.LogObjectAsJson(span, "request", request)

	unlinkEmailEvent, err := event.NewContactUnlinkEmailEvent(a, request.Email)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "NewContactUnlinkEmailEvent")
	}
	eventstore.EnrichEventWithMetadataExtended(&unlinkEmailEvent, span, eventstore.EventMetadata{
		Tenant: a.GetTenant(),
		UserId: request.LoggedInUserId,
		App:    request.AppSource,
	})

	return a.Apply(unlinkEmailEvent)
}

func (a *ContactAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContactCreateV1:
		return a.onContactCreate(evt)
	case event.ContactUpdateV1:
		return a.onContactUpdate(evt)
	case event.ContactPhoneNumberLinkV1:
		return a.onPhoneNumberLink(evt)
	case event.ContactEmailLinkV1:
		return a.onEmailLink(evt)
	case event.ContactEmailUnlinkV1:
		return a.onEmailUnlink(evt)
	case event.ContactLocationLinkV1:
		return a.onLocationLink(evt)
	case event.ContactOrganizationLinkV1:
		return a.onOrganizationLink(evt)
	case event.ContactAddSocialV1:
		return a.onAddSocial(evt)
	case event.ContactRemoveSocialV1:
		return a.onRemoveSocial(evt)
	case event.ContactAddTagV1:
		return a.onContactAddTag(evt)
	case event.ContactRemoveTagV1:
		return a.onContactRemoveTag(evt)
	case event.ContactAddLocationV1:
		return a.onAddLocation(evt)
	default:
		return nil
	}
}

func (a *ContactAggregate) onContactCreate(evt eventstore.Event) error {
	var eventData event.ContactCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.FirstName = eventData.FirstName
	a.Contact.LastName = eventData.LastName
	a.Contact.Prefix = eventData.Prefix
	a.Contact.Name = eventData.Name
	a.Contact.Description = eventData.Description
	a.Contact.Timezone = eventData.Timezone
	a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
	a.Contact.Username = eventData.Username
	a.Contact.Source = cmnmod.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Contact.CreatedAt = eventData.CreatedAt
	a.Contact.UpdatedAt = eventData.UpdatedAt
	if eventData.ExternalSystem.Available() {
		a.Contact.ExternalSystems = []cmnmod.ExternalSystem{eventData.ExternalSystem}
	}
	return nil
}

func (a *ContactAggregate) onContactUpdate(evt eventstore.Event) error {
	var eventData event.ContactUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	if eventData.Source != a.Contact.Source.SourceOfTruth && a.Contact.Source.SourceOfTruth == constants.SourceOpenline {
		if a.Contact.Name == "" {
			a.Contact.Name = eventData.Name
		}
		if a.Contact.FirstName == "" {
			a.Contact.FirstName = eventData.FirstName
		}
		if a.Contact.LastName == "" {
			a.Contact.LastName = eventData.LastName
		}
		if a.Contact.Timezone == "" {
			a.Contact.Timezone = eventData.Timezone
		}
		if a.Contact.ProfilePhotoUrl == "" {
			a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
		}
		if a.Contact.Username == "" {
			a.Contact.Username = eventData.Username
		}
		if a.Contact.Prefix == "" {
			a.Contact.Prefix = eventData.Prefix
		}
	} else {
		a.Contact.Name = eventData.Name
		a.Contact.FirstName = eventData.FirstName
		a.Contact.LastName = eventData.LastName
		a.Contact.Prefix = eventData.Prefix
		a.Contact.Description = eventData.Description
		a.Contact.Timezone = eventData.Timezone
		a.Contact.ProfilePhotoUrl = eventData.ProfilePhotoUrl
		a.Contact.Username = eventData.Username
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	if eventData.Source == constants.SourceOpenline {
		a.Contact.Source.SourceOfTruth = eventData.Source
	}

	if eventData.ExternalSystem.Available() {
		found := false
		for _, externalSystem := range a.Contact.ExternalSystems {
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
			a.Contact.ExternalSystems = append(a.Contact.ExternalSystems, eventData.ExternalSystem)
		}
	}

	return nil
}

func (a *ContactAggregate) onPhoneNumberLink(evt eventstore.Event) error {
	var eventData event.ContactLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.PhoneNumbers == nil {
		a.Contact.PhoneNumbers = make(map[string]ContactPhoneNumber)
	}
	a.Contact.PhoneNumbers[eventData.PhoneNumberId] = ContactPhoneNumber{
		Label:   eventData.Label,
		Primary: eventData.Primary,
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onEmailLink(evt eventstore.Event) error {
	var eventData event.ContactLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.Emails == nil {
		a.Contact.Emails = make(map[string]ContactEmail)
	}
	a.Contact.Emails[eventData.EmailId] = ContactEmail{
		Primary: eventData.Primary,
	}
	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onEmailUnlink(evt eventstore.Event) error {
	var eventData event.ContactUnlinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.Emails = make(map[string]ContactEmail)
	return nil
}

func (a *ContactAggregate) onLocationLink(evt eventstore.Event) error {
	var eventData event.ContactLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Contact.LocationIds = utils.AddToListIfNotExists(a.Contact.LocationIds, eventData.LocationId)
	return nil
}

func (a *ContactAggregate) onOrganizationLink(evt eventstore.Event) error {
	var eventData event.ContactLinkWithOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.JobRolesByOrganization == nil {
		a.Contact.JobRolesByOrganization = make(map[string]JobRole)
	}
	jobRole, found := a.Contact.JobRolesByOrganization[eventData.OrganizationId]
	if !found {
		a.Contact.JobRolesByOrganization[eventData.OrganizationId] = JobRole{
			JobTitle:    eventData.JobTitle,
			Primary:     eventData.Primary,
			Description: eventData.Description,
			StartedAt:   eventData.StartedAt,
			EndedAt:     eventData.EndedAt,
			Source: cmnmod.Source{
				Source:        eventData.SourceFields.Source,
				SourceOfTruth: eventData.SourceFields.SourceOfTruth,
				AppSource:     eventData.SourceFields.AppSource,
			},
			CreatedAt: eventData.CreatedAt,
		}
	} else {
		if eventData.SourceFields.Source != jobRole.Source.SourceOfTruth && jobRole.Source.SourceOfTruth == constants.SourceOpenline {
			if jobRole.JobTitle == "" {
				jobRole.JobTitle = eventData.JobTitle
			}
			if jobRole.Description == "" {
				jobRole.Description = eventData.Description
			}
			if jobRole.StartedAt == nil {
				jobRole.StartedAt = eventData.StartedAt
			}
			if jobRole.EndedAt == nil {
				jobRole.EndedAt = eventData.EndedAt
			}
		} else {
			jobRole.JobTitle = eventData.JobTitle
			jobRole.Primary = eventData.Primary
			jobRole.Description = eventData.Description
			jobRole.StartedAt = eventData.StartedAt
			jobRole.EndedAt = eventData.EndedAt
		}
	}

	a.Contact.UpdatedAt = eventData.UpdatedAt
	return nil
}

func (a *ContactAggregate) onAddSocial(evt eventstore.Event) error {
	var eventData event.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.Socials == nil {
		a.Contact.Socials = make(map[string]cmnmod.Social)
	}
	a.Contact.Socials[eventData.SocialId] = cmnmod.Social{
		Url:            eventData.Url,
		Alias:          eventData.Alias,
		ExternalId:     eventData.ExternalId,
		FollowersCount: eventData.FollowersCount,
	}
	return nil
}

func (a *ContactAggregate) onRemoveSocial(evt eventstore.Event) error {
	var eventData event.ContactAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.Socials == nil {
		a.Contact.Socials = make(map[string]cmnmod.Social)
	}
	delete(a.Contact.Socials, eventData.SocialId)
	return nil
}

func (a *ContactAggregate) onContactAddTag(evt eventstore.Event) error {
	var eventData event.ContactAddTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contact.TagIds = append(a.Contact.TagIds, eventData.TagId)
	a.Contact.TagIds = utils.RemoveDuplicates(a.Contact.TagIds)

	return nil
}

func (a *ContactAggregate) onContactRemoveTag(evt eventstore.Event) error {
	var eventData event.ContactRemoveTagEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}

	a.Contact.TagIds = utils.RemoveFromList(a.Contact.TagIds, eventData.TagId)

	return nil
}

func (a *ContactAggregate) onAddLocation(evt eventstore.Event) error {
	var eventData event.ContactAddLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	if a.Contact.Locations == nil {
		a.Contact.Locations = make(map[string]cmnmod.Location)
	}
	a.Contact.Locations[eventData.LocationId] = cmnmod.Location{
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
