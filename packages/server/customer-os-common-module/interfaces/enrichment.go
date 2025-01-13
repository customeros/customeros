package interfaces

import "context"

type EnrichmentService interface {
	Snitcher(ctx context.Context, ipAddress string) (*SnitcherDataResponse, error)
}

type SnitcherDataResponse struct {
	Status       string          `json:"status"`
	Message      string          `json:"message,omitempty"`
	CompanyFound bool            `json:"companyFound"`
	Data         *CompanyDetails `json:"data,omitempty"`
}

type CompanyDetails struct {
	Name          string          `json:"name"`
	Domain        string          `json:"domain"`
	Website       string          `json:"website"`
	Industry      string          `json:"industry"`
	FoundedYear   interface{}     `json:"founded_year"`
	EmployeeRange string          `json:"employee_range"`
	AnnualRevenue interface{}     `json:"annual_revenue"`
	TotalFunding  interface{}     `json:"total_funding"`
	Location      Location        `json:"location"`
	Description   string          `json:"description"`
	Phone         string          `json:"phone"`
	Geo           GeoLocation     `json:"geo"`
	Profiles      *SocialProfiles `json:"profiles"`
}

type GeoLocation struct {
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	State        string  `json:"state"`
	StateCode    *string `json:"state_code"`
	PostalCode   *string `json:"postal_code"`
	City         string  `json:"city"`
	Street       *string `json:"street"`
	StreetNumber *string `json:"street_number"`
}

type Location struct {
	CityName     string `json:"cityName"`
	RegionName   string `json:"regionName"`
	PostalCode   string `json:"postalCode"`
	StreetName   string `json:"streetName"`
	StreetNumber string `json:"streetNumber"`
	Country      struct {
		Name string `json:"name"`
		Iso2 string `json:"iso2"`
		Iso3 string `json:"iso3"`
	} `json:"country"`
	RawLocation string
}

type SocialProfiles struct {
	Crunchbase *ProfileInfo `json:"crunchbase"`
	LinkedIn   *ProfileInfo `json:"linkedin"`
	Facebook   *ProfileInfo `json:"facebook"`
}

// ProfileInfo represents common social media profile attributes
type ProfileInfo struct {
	Handle string      `json:"handle"`
	URL    interface{} `json:"url"`
}
