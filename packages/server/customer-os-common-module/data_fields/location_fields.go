package data_fields

import "time"

type LocationFields struct {
	Source        *string    `json:"source,omitempty"`
	AppSource     *string    `json:"appSource,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	RawAddress    string     `json:"rawAddress,omitempty"`
	Name          string     `json:"name,omitempty"`
	Country       string     `json:"country,omitempty"`
	CountryCodeA2 string     `json:"countryCodeA2,omitempty"`
	CountryCodeA3 string     `json:"countryCodeA3,omitempty"`
	Region        string     `json:"region,omitempty"`
	Locality      string     `json:"locality,omitempty"`
	Address       string     `json:"address,omitempty"`
	Address2      string     `json:"address2"`
	Zip           string     `json:"zip"`
	AddressType   string     `json:"addressType"`
	HouseNumber   string     `json:"houseNumber"`
	PostalCode    string     `json:"postalCode"`
	PlusFour      string     `json:"plusFour"`
	Commercial    bool       `json:"commercial"`
	Predirection  string     `json:"predirection"`
	District      string     `json:"district"`
	Street        string     `json:"street"`
	Latitude      *float64   `json:"latitude"`
	Longitude     *float64   `json:"longitude"`
	TimeZone      string     `json:"timeZone"`
	UtcOffset     *float64   `json:"utcOffset"`
}
