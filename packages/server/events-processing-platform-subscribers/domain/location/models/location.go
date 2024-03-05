package models

import (
	"fmt"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"time"
)

type LocationValidation struct {
	ValidationError string `json:"validationError"`
	SkipReason      string `json:"skipReason"`
}

type LocationAddress struct {
	Country      string   `json:"country"`
	Region       string   `json:"region"`
	District     string   `json:"district"`
	Locality     string   `json:"locality"`
	Street       string   `json:"street"`
	Address1     string   `json:"address1"`
	Address2     string   `json:"address2"`
	Zip          string   `json:"zip"`
	AddressType  string   `json:"addressType"`
	HouseNumber  string   `json:"houseNumber"`
	PostalCode   string   `json:"postalCode"`
	PlusFour     string   `json:"plusFour"`
	Commercial   bool     `json:"commercial"`
	Predirection string   `json:"predirection"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	TimeZone     string   `json:"timeZone"`
	UtcOffset    int      `json:"utcOffset"`
}

func (l *LocationAddress) FillFrom(newLoc LocationAddress) {
	if l.Country == "" && newLoc.Country != "" {
		l.Country = newLoc.Country
	}
	if l.Region == "" && newLoc.Region != "" {
		l.Region = newLoc.Region
	}
	if l.District == "" && newLoc.District != "" {
		l.District = newLoc.District
	}
	if l.Locality == "" && newLoc.Locality != "" {
		l.Locality = newLoc.Locality
	}
	if l.Street == "" && newLoc.Street != "" {
		l.Street = newLoc.Street
	}
	if l.Address1 == "" && newLoc.Address1 != "" {
		l.Address1 = newLoc.Address1
	}
	if l.Address2 == "" && newLoc.Address2 != "" {
		l.Address2 = newLoc.Address2
	}
	if l.Zip == "" && newLoc.Zip != "" {
		l.Zip = newLoc.Zip
	}
	if l.AddressType == "" && newLoc.AddressType != "" {
		l.AddressType = newLoc.AddressType
	}
	if l.HouseNumber == "" && newLoc.HouseNumber != "" {
		l.HouseNumber = newLoc.HouseNumber
	}
	if l.PlusFour == "" && newLoc.PlusFour != "" {
		l.PlusFour = newLoc.PlusFour
	}
	if l.PostalCode == "" && newLoc.PostalCode != "" {
		l.PostalCode = newLoc.PostalCode
	}
	if l.Predirection == "" && newLoc.Predirection != "" {
		l.Predirection = newLoc.Predirection
	}
	if l.Latitude == nil && newLoc.Latitude != nil {
		l.Latitude = newLoc.Latitude
	}
	if l.Longitude == nil && newLoc.Longitude != nil {
		l.Longitude = newLoc.Longitude
	}
	if l.TimeZone == "" && newLoc.TimeZone != "" {
		l.TimeZone = newLoc.TimeZone
	}
	l.UtcOffset = newLoc.UtcOffset
	l.Commercial = newLoc.Commercial
}

func (l *LocationAddress) From(fields LocationAddressFields) {
	l.Country = fields.Country
	l.Region = fields.Region
	l.District = fields.District
	l.Locality = fields.Locality
	l.Street = fields.Street
	l.Address1 = fields.Address1
	l.Address2 = fields.Address2
	l.Zip = fields.Zip
	l.AddressType = fields.AddressType
	l.HouseNumber = fields.HouseNumber
	l.PostalCode = fields.PostalCode
	l.PlusFour = fields.PlusFour
	l.Commercial = fields.Commercial
	l.Predirection = fields.Predirection
	l.Latitude = fields.Latitude
	l.Longitude = fields.Longitude
	l.TimeZone = fields.TimeZone
	l.UtcOffset = fields.UtcOffset
}

type Location struct {
	ID                 string             `json:"id"`
	Source             cmnmod.Source      `json:"source"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	LocationValidation LocationValidation `json:"locationValidation"`
	Name               string             `json:"name"`
	RawAddress         string             `json:"rawAddress"`
	LocationAddress    LocationAddress    `json:"locationAddress"`
}

func (l *Location) String() string {
	return fmt.Sprintf("ID: %s, Source: %s, CreatedAt: %s, UpdatedAt: %s, LocationValidation: %s, Name: %s, RawAddress: %s, LocationAddress: %v", l.ID, l.Source, l.CreatedAt, l.UpdatedAt, l.LocationValidation, l.Name, l.RawAddress, l.LocationAddress)
}

func (l *LocationValidation) String() string {
	return fmt.Sprintf("ValidationError: %s, SkipReason: %s", l.ValidationError, l.SkipReason)
}

func (l *LocationAddress) String() string {
	return fmt.Sprintf("Country: %s, Region: %s, District: %s, Locality: %s, Street: %s, Address1: %s, Address2: %s, Zip: %s, AddressType: %s, HouseNumber: %s, PostalCode: %s, PlusFour: %s, Commercial: %t, Predirection: %s, Latitude: %v, Longitude: %v, TimeZone: %s, UtcOffset: %d", l.Country, l.Region, l.District, l.Locality, l.Street, l.Address1, l.Address2, l.Zip, l.AddressType, l.HouseNumber, l.PostalCode, l.PlusFour, l.Commercial, l.Predirection, l.Latitude, l.Longitude, l.TimeZone, l.UtcOffset)
}

func NewLocation() *Location {
	return &Location{}
}
