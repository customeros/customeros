package graph

import (
	"context"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/service"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/aggregate"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type LocationEventHandler struct {
	services *service.Services
}

func NewLocationEventHandler(services *service.Services) *LocationEventHandler {
	return &LocationEventHandler{
		services: services,
	}
}

func (h *LocationEventHandler) OnLocationValidated(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationValidated")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationValidatedEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	locationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	data := locationAddressToAddressDetails(eventData.LocationAddress)
	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.LocationValidated(ctx, eventData.Tenant, locationId, data, eventData.ValidatedAt)

	return err
}

func (h *LocationEventHandler) OnLocationValidationFailed(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.OnLocationValidationFailed")
	defer span.Finish()
	setEventSpanTagsAndLogFields(span, evt)

	var eventData events.LocationFailedValidationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	locationId := aggregate.GetLocationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.services.CommonServices.Neo4jRepositories.LocationWriteRepository.FailLocationValidation(ctx, eventData.Tenant, locationId, eventData.ValidationError, eventData.ValidatedAt)

	return err
}

func locationAddressToAddressDetails(locationAddress models.LocationAddress) neo4jrepository.AddressDetails {
	return neo4jrepository.AddressDetails{
		Latitude:     locationAddress.Latitude,
		Longitude:    locationAddress.Longitude,
		Country:      locationAddress.Country,
		Region:       locationAddress.Region,
		District:     locationAddress.District,
		Locality:     locationAddress.Locality,
		Street:       locationAddress.Street,
		Address:      locationAddress.Address1,
		Address2:     locationAddress.Address2,
		Zip:          locationAddress.Zip,
		AddressType:  locationAddress.AddressType,
		HouseNumber:  locationAddress.HouseNumber,
		PostalCode:   locationAddress.PostalCode,
		PlusFour:     locationAddress.PlusFour,
		Commercial:   locationAddress.Commercial,
		Predirection: locationAddress.Predirection,
		TimeZone:     locationAddress.TimeZone,
		UtcOffset:    locationAddress.UtcOffset,
	}

}
