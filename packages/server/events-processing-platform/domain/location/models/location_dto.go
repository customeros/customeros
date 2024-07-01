package models

type LocationAddressFields struct {
	Country      string
	Region       string
	District     string
	Locality     string
	Street       string
	Address1     string
	Address2     string
	Zip          string
	AddressType  string
	HouseNumber  string
	PostalCode   string
	PlusFour     string
	Commercial   bool
	Predirection string
	Latitude     *float64
	Longitude    *float64
	TimeZone     string
	UtcOffset    *float64
}
