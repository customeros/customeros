package entity

import (
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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

	SourceOfTruth neo4jentity.DataSource
	Source        neo4jentity.DataSource
	AppSource     string
}

type LocationEntities []LocationEntity
