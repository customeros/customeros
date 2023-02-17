package entity

import (
	"fmt"
	"time"
)

type LocationEntity struct {
	Id            string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Country       string
	Region        string
	Locality      string
	Address       string
	Address2      string
	Zip           string
	AddressType   string
	HouseNumber   string
	PostalCode    string
	PlusFour      string
	Commercial    bool
	Predirection  string
	District      string
	Street        string
	RawAddress    string
	Latitude      *float64
	Longitude     *float64
	SourceOfTruth DataSource
	Source        DataSource
	AppSource     string

	DataloaderKey string
}

func (location LocationEntity) ToString() string {
	return fmt.Sprintf("id: %s name: %s", location.Id, location.Name)
}

type LocationEntities []LocationEntity

func (location LocationEntity) Labels(tenant string) []string {
	return []string{"Location", "Location_" + tenant}
}
