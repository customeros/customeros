package location

type Location struct {
	Country       string   `json:"country"`
	CountryCodeA2 string   `json:"countryCodeA2"`
	CountryCodeA3 string   `json:"countryCodeA3"`
	Region        string   `json:"region"`
	Locality      string   `json:"locality"`
	Address       string   `json:"address"`
	Address2      string   `json:"address2"`
	Zip           string   `json:"zip"`
	AddressType   string   `json:"addressType"`
	HouseNumber   string   `json:"houseNumber"`
	PostalCode    string   `json:"postalCode"`
	PlusFour      string   `json:"plusFour"`
	Commercial    bool     `json:"commercial"`
	Predirection  string   `json:"predirection"`
	District      string   `json:"district"`
	Street        string   `json:"street"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	TimeZone      string   `json:"timeZone"`
	UtcOffset     *float64 `json:"utcOffset"`
}
