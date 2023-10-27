package entity

import (
	"time"
)

type LocationEntity struct {
	Id           string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Country      string `neo4jDb:"property:country;lookupName:COUNTRY;supportCaseSensitive:true"`
	Region       string `neo4jDb:"property:region;lookupName:REGION;supportCaseSensitive:true"`
	Locality     string `neo4jDb:"property:locality;lookupName:LOCALITY;supportCaseSensitive:true"`
	Address      string
	Address2     string
	Zip          string
	AddressType  string
	HouseNumber  string
	PostalCode   string
	PlusFour     string
	Commercial   bool
	Predirection string
	District     string
	Street       string
	RawAddress   string
	Latitude     *float64
	Longitude    *float64
	TimeZone     string
	UtcOffset    int64

	SourceOfTruth DataSource
	Source        DataSource
	AppSource     string
}

type LocationEntities []LocationEntity
