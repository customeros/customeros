package events

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type OrganizationAddLocationEvent struct {
	Tenant        string    `json:"tenant" validate:"required"`
	LocationId    string    `json:"locationId" validate:"required"`
	Source        string    `json:"source"`
	SourceOfTruth string    `json:"sourceOfTruth"`
	AppSource     string    `json:"appSource"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	RawAddress    string    `json:"rawAddress"`
	Country       string    `json:"country"`
	CountryCodeA2 string    `json:"countryCodeA2"`
	CountryCodeA3 string    `json:"countryCodeA3"`
	Region        string    `json:"region"`
	District      string    `json:"district"`
	Locality      string    `json:"locality"`
	AddressLine1  string    `json:"addressLine1"`
	AddressLine2  string    `json:"addressLine2"`
	Street        string    `json:"street"`
	HouseNumber   string    `json:"houseNumber"`
	ZipCode       string    `json:"zipCode"`
	PostalCode    string    `json:"postalCode"`
	AddressType   string    `json:"addressType"`
	Commercial    bool      `json:"commercial"`
	Predirection  string    `json:"predirection"`
	PlusFour      string    `json:"plusFour"`
	TimeZone      string    `json:"timeZone"`
	UtcOffset     *float64  `json:"utcOffset"`
	Latitude      *float64  `json:"latitude"`
	Longitude     *float64  `json:"longitude"`
}

func NewOrganizationAddLocationEvent(aggregate eventstore.Aggregate, locationId string, location cmnmod.Location, sourceFields events.Source, createdAt time.Time) (eventstore.Event, error) {
	eventData := OrganizationAddLocationEvent{
		Tenant:        aggregate.GetTenant(),
		LocationId:    locationId,
		Source:        sourceFields.Source,
		SourceOfTruth: sourceFields.SourceOfTruth,
		AppSource:     sourceFields.AppSource,
		CreatedAt:     createdAt,
		Name:          location.Name,
		RawAddress:    location.RawAddress,
		Country:       location.Country,
		CountryCodeA2: location.CountryCodeA2,
		CountryCodeA3: location.CountryCodeA3,
		Region:        location.Region,
		District:      location.District,
		Locality:      location.Locality,
		AddressLine1:  location.AddressLine1,
		AddressLine2:  location.AddressLine2,
		Street:        location.Street,
		HouseNumber:   location.HouseNumber,
		ZipCode:       location.ZipCode,
		PostalCode:    location.PostalCode,
		AddressType:   location.AddressType,
		Commercial:    location.Commercial,
		Predirection:  location.Predirection,
		PlusFour:      location.PlusFour,
		TimeZone:      location.TimeZone,
		UtcOffset:     location.UtcOffset,
		Latitude:      location.Latitude,
		Longitude:     location.Longitude,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate OrganizationAddLocationEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, OrganizationAddLocationV1)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for OrganizationAddLocationEvent")
	}
	return event, nil
}
