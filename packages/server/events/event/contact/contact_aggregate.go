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
	ContactAggregateType eventstore.AggregateType = "contact"
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
	case *contactpb.ContactRemoveSocialGrpcRequest:
		return nil, a.removeSocial(ctx, r)
	case *contactpb.ContactAddLocationGrpcRequest:
		return a.addLocation(ctx, r)
	default:
		tracing.TraceErr(span, eventstore.ErrInvalidRequestType)
		return nil, eventstore.ErrInvalidRequestType
	}
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
		return errors.New("SocialId is required")
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

func (a *ContactAggregate) When(evt eventstore.Event) error {
	switch evt.GetEventType() {
	case event.ContactAddTagV1,
		event.ContactRemoveTagV1,
		event.ContactEmailLinkV1,
		event.ContactEmailUnlinkV1,
		event.ContactUpdateV1,
		event.ContactCreateV1:
		return nil

	case event.ContactPhoneNumberLinkV1:
		return a.onPhoneNumberLink(evt)
	case event.ContactLocationLinkV1:
		return a.onLocationLink(evt)
	case event.ContactOrganizationLinkV1:
		return a.onOrganizationLink(evt)
	case event.ContactAddLocationV1:
		return a.onAddLocation(evt)
	default:
		return nil
	}
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
