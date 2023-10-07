package entity

import (
	"time"
)

type LocationEntity struct {
	Id           string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Country      string
	Region       string
	Locality     string
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
