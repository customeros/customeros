package common

type Location struct {
	Name          string   `json:"name"`
	RawAddress    string   `json:"rawAddress"`
	Country       string   `json:"country"`
	CountryCodeA2 string   `json:"countryCodeA2"`
	CountryCodeA3 string   `json:"countryCodeA3"`
	Region        string   `json:"region"`
	District      string   `json:"district"`
	Locality      string   `json:"locality"`
	AddressLine1  string   `json:"addressLine1"`
	AddressLine2  string   `json:"addressLine2"`
	Street        string   `json:"street"`
	HouseNumber   string   `json:"houseNumber"`
	ZipCode       string   `json:"zipCode"`
	PostalCode    string   `json:"postalCode"`
	AddressType   string   `json:"addressType"`
	Commercial    bool     `json:"commercial"`
	Predirection  string   `json:"predirection"`
	PlusFour      string   `json:"plusFour"`
	TimeZone      string   `json:"timeZone"`
	UtcOffset     *float64 `json:"utcOffset"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
}

func (l Location) IsEmpty() bool {
	return l.Name == "" && l.RawAddress == "" && l.Country == "" && l.CountryCodeA2 == "" && l.CountryCodeA3 == "" &&
		l.Region == "" && l.District == "" && l.Locality == "" && l.AddressLine1 == "" && l.AddressLine2 == "" &&
		l.Street == "" && l.HouseNumber == "" && l.ZipCode == "" && l.PostalCode == "" && l.AddressType == "" &&
		l.Predirection == "" && l.PlusFour == "" && l.TimeZone == "" && l.UtcOffset == nil && l.Latitude == nil && l.Longitude == nil
}
