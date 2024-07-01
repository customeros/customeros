package entity

import (
	"time"
)

type LocationProperty string

const (
	LocationPropertyName          LocationProperty = "name"
	LocationPropertyCountry       LocationProperty = "country"
	LocationPropertyRegion        LocationProperty = "region"
	LocationPropertyLocality      LocationProperty = "locality"
	LocationPropertyAddress       LocationProperty = "address"
	LocationPropertyAddress2      LocationProperty = "address2"
	LocationPropertyZip           LocationProperty = "zip"
	LocationPropertyAddressType   LocationProperty = "addressType"
	LocationPropertyHouseNumber   LocationProperty = "houseNumber"
	LocationPropertyPostalCode    LocationProperty = "postalCode"
	LocationPropertyPlusFour      LocationProperty = "plusFour"
	LocationPropertyCommercial    LocationProperty = "commercial"
	LocationPropertyPredirection  LocationProperty = "predirection"
	LocationPropertyDistrict      LocationProperty = "district"
	LocationPropertyStreet        LocationProperty = "street"
	LocationPropertyRawAddress    LocationProperty = "rawAddress"
	LocationPropertyLatitude      LocationProperty = "latitude"
	LocationPropertyLongitude     LocationProperty = "longitude"
	LocationPropertyTimeZone      LocationProperty = "timeZone"
	LocationPropertyUtcOffset     LocationProperty = "utcOffset"
	LocationPropertyCountryCodeA2 LocationProperty = "countryCodeA2"
	LocationPropertyCountryCodeA3 LocationProperty = "countryCodeA3"
)

type LocationEntity struct {
	DataLoaderKey

	Id            string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Country       string `neo4jDb:"property:country;lookupName:COUNTRY;supportCaseSensitive:true"`
	CountryCodeA2 string `neo4jDb:"property:countryCodeA2;lookupName:COUNTRY_CODE_A2;supportCaseSensitive:true"`
	CountryCodeA3 string `neo4jDb:"property:countryCodeA3;lookupName:COUNTRY_CODE_A3;supportCaseSensitive:true"`
	Region        string `neo4jDb:"property:region;lookupName:REGION;supportCaseSensitive:true"`
	Locality      string `neo4jDb:"property:locality;lookupName:LOCALITY;supportCaseSensitive:true"`
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
	TimeZone      string
	UtcOffset     *float64
	SourceOfTruth DataSource
	Source        DataSource
	AppSource     string
}

type LocationEntities []LocationEntity
