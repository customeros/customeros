package aggregate

import (
	commonModels "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/pkg/errors"
)

const (
	LocationAggregateType eventstore.AggregateType = "location"
)

type LocationAggregate struct {
	*eventstore.AggregateBase
	Location *models.Location
}

func NewLocationAggregateWithTenantAndID(tenant, id string) *LocationAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewLocationAggregate()
	aggregate.SetID(tenant + "-" + id)
	return aggregate
}

func NewLocationAggregate() *LocationAggregate {
	locationAggregate := &LocationAggregate{Location: models.NewLocation()}
	base := eventstore.NewAggregateBase(locationAggregate.When)
	base.SetType(LocationAggregateType)
	locationAggregate.AggregateBase = base
	return locationAggregate
}

func (a *LocationAggregate) When(event eventstore.Event) error {
	switch event.GetEventType() {
	case
		events.LocationCreateV1Legacy,
		events.LocationCreateV1:
		return a.onLocationCreate(event)
	case
		events.LocationUpdateV1Legacy,
		events.LocationUpdateV1:
		return a.onLocationUpdate(event)
	case events.LocationValidationSkippedV1:
		return a.OnLocationSkippedValidation(event)
	case events.LocationValidationFailedV1:
		return a.OnLocationFailedValidation(event)
	case events.LocationValidatedV1:
		return a.OnLocationValidated(event)

	default:
		err := eventstore.ErrInvalidEventType
		err.EventType = event.GetEventType()
		return err
	}
}

func (a *LocationAggregate) onLocationCreate(event eventstore.Event) error {
	var eventData events.LocationCreateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Location.Source = commonModels.Source{
		Source:        eventData.Source,
		SourceOfTruth: eventData.SourceOfTruth,
		AppSource:     eventData.AppSource,
	}
	a.Location.CreatedAt = eventData.CreatedAt
	a.Location.UpdatedAt = eventData.UpdatedAt
	a.Location.Name = eventData.Name
	a.Location.RawAddress = eventData.RawAddress
	a.Location.LocationAddress = eventData.LocationAddress
	return nil
}

func (a *LocationAggregate) onLocationUpdate(event eventstore.Event) error {
	var eventData events.LocationUpdateEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Location.Source.SourceOfTruth = eventData.SourceOfTruth
	a.Location.UpdatedAt = eventData.UpdatedAt
	a.Location.Name = eventData.Name
	a.Location.RawAddress = eventData.RawAddress
	a.Location.LocationAddress.Country = eventData.LocationAddress.Country
	a.Location.LocationAddress.Region = eventData.LocationAddress.Region
	a.Location.LocationAddress.District = eventData.LocationAddress.District
	a.Location.LocationAddress.Locality = eventData.LocationAddress.Locality
	a.Location.LocationAddress.Street = eventData.LocationAddress.Street
	a.Location.LocationAddress.Address1 = eventData.LocationAddress.Address1
	a.Location.LocationAddress.Address2 = eventData.LocationAddress.Address2
	a.Location.LocationAddress.Zip = eventData.LocationAddress.Zip
	a.Location.LocationAddress.AddressType = eventData.LocationAddress.AddressType
	a.Location.LocationAddress.HouseNumber = eventData.LocationAddress.HouseNumber
	a.Location.LocationAddress.PlusFour = eventData.LocationAddress.PlusFour
	a.Location.LocationAddress.PostalCode = eventData.LocationAddress.PostalCode
	a.Location.LocationAddress.Commercial = eventData.LocationAddress.Commercial
	a.Location.LocationAddress.Predirection = eventData.LocationAddress.Predirection
	a.Location.LocationAddress.Latitude = eventData.LocationAddress.Latitude
	a.Location.LocationAddress.Longitude = eventData.LocationAddress.Longitude

	return nil
}

func (a *LocationAggregate) OnLocationSkippedValidation(event eventstore.Event) error {
	var eventData events.LocationSkippedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Location.LocationValidation.SkipReason = eventData.Reason
	a.Location.UpdatedAt = eventData.ValidatedAt
	return nil
}

func (a *LocationAggregate) OnLocationFailedValidation(event eventstore.Event) error {
	var eventData events.LocationFailedValidationEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Location.LocationValidation.ValidationError = eventData.ValidationError
	a.Location.UpdatedAt = eventData.ValidatedAt
	return nil
}

func (a *LocationAggregate) OnLocationValidated(event eventstore.Event) error {
	var eventData events.LocationValidatedEvent
	if err := event.GetJsonData(&eventData); err != nil {
		return errors.Wrap(err, "GetJsonData")
	}
	a.Location.LocationValidation.ValidationError = ""
	a.Location.LocationValidation.SkipReason = ""
	a.Location.UpdatedAt = eventData.ValidatedAt
	a.Location.LocationAddress.FillFrom(eventData.LocationAddress)
	return nil
}
