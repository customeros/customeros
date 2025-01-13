package interfaces

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type EnrichmentService interface {
	//primary interfaces
	EnrichOrganization(ctx context.Context, domain, linkedinURL *string) (*OrganizationData, error)
	EnrichPerson(ctx context.Context, person PersonSearch) (*uint64, *entity.ScrapInResponseBody, error)
	IPIdentity(ctx context.Context, ipAddress string) (*SnitcherResponse, error)
	FindWorkEmail(ctx context.Context, linkedInUrl, firstName, lastName, companyName, companyDomain string, enrichPhoneNumber bool) (dbID string, betterContactRequstID string, response *entity.BetterContactResponseBody, err error)

	// only use if you must
	GetBrandfetchByDomain(ctx context.Context, domain string) (*entity.BrandfetchResponseBody, error)

	ScrapInPersonProfile(ctx context.Context, linkedInUrl string) (uint64, *entity.ScrapInResponseBody, error)
	ScrapInSearchPerson(ctx context.Context, email, fistName, lastName, domain, companyName string) (uint64, *entity.ScrapInResponseBody, error)
	ScrapInCompanyProfile(ctx context.Context, linkedInUrl string) (uint64, *entity.ScrapInResponseBody, error)
	ScrapInSearchCompany(ctx context.Context, domain string) (uint64, *entity.ScrapInResponseBody, error)
}

type PersonSearch struct {
	LinkedinURL *string
	FirstName   *string
	LastName    *string
	Email       *string
	Domain      *string
	CompanyName *string
}

type OrganizationData struct {
	Name             string               `json:"name"`
	Domain           string               `json:"domain"`
	ShortDescription string               `json:"description"`
	LongDescription  string               `json:"longDescription"`
	Website          string               `json:"website"`
	Employees        int64                `json:"employees"`
	FoundedYear      int64                `json:"foundedYear"`
	Public           *bool                `json:"public,omitempty"`
	Logos            []string             `json:"logos"`
	Icons            []string             `json:"icons"`
	Industry         string               `json:"industry"`
	Socials          []OrganizationSocial `json:"socials"`
	Location         OrganizationLocation `json:"location"`
}

type OrganizationSocial struct {
	Url   string `json:"url"`
	Alias string `json:"alias"`
	Id    string `json:"id"`
}

type OrganizationLocation struct {
	IsHeadquarter *bool  `json:"isHeadquarter"`
	Country       string `json:"country"`
	CountryCodeA2 string `json:"countryCodeA2"`
	CountryCodeA3 string `json:"countryCodeA3"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postalCode"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
}

func (l OrganizationLocation) IsEmpty() bool {
	return l.Country == "" && l.CountryCodeA2 == "" && l.Locality == "" && l.Region == ""
}

// SnitcherResponse represents the top-level response structure
type SnitcherResponse struct {
	Fuzzy   bool            `json:"fuzzy"`
	Domain  string          `json:"domain"`
	Type    string          `json:"type"`
	Company *CompanyDetails `json:"company"`
}

// CompanyDetails contains the main company information
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

// GeoLocation represents geographical information
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
	RawLocation string // To store the raw string value if the location is a single string
}

// SocialProfiles contains social media profile information
type SocialProfiles struct {
	Crunchbase *ProfileInfo `json:"crunchbase"`
	LinkedIn   *ProfileInfo `json:"linkedin"`
	Facebook   *ProfileInfo `json:"facebook"`
}

// ProfileInfo represents common social media profile attributes
type ProfileInfo struct {
	Handle string      `json:"handle"`
	URL    interface{} `json:"url"` // Consider using *string if possible
}

// Implement the UnmarshalJSON method for Location
func (l *Location) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal into the struct form first
	type Alias Location
	var alias Alias
	if err := json.Unmarshal(data, &alias); err == nil {
		// If successful, copy the unmarshaled data to the current object
		*l = Location(alias)
		return nil
	}

	// If unmarshaling to struct fails, assume it's a raw string
	var rawString string
	if err := json.Unmarshal(data, &rawString); err == nil {
		l.RawLocation = rawString
		return nil
	}

	// If both attempts fail, return an error
	return fmt.Errorf("invalid location format")
}

func (s *SnitcherResponse) CompanyName() string {
	if s.Company != nil && s.Company.Name != "" {
		return s.Company.Name
	}
	return ""
}

func (s *SnitcherResponse) CompanyDomain() string {
	if s.Company != nil && s.Company.Domain != "" {
		return s.Company.Domain
	}
	return ""
}

func (s *SnitcherResponse) LinkedinSlug() string {
	if s.Company != nil && s.Company.Profiles.LinkedIn.Handle != "" {
		return s.Company.Profiles.LinkedIn.Handle
	}
	return ""
}

func (s *SnitcherResponse) CompanyWebsite() string {
	if s.Company != nil && s.Company.Website != "" {
		return s.Company.Website
	}
	return ""
}

func (s *SnitcherResponse) CompanyFound() bool {
	if s.Company == nil {
		return false
	}
	return true
}

type SnitcherDataResponse struct {
	Status       string          `json:"status"`
	Message      string          `json:"message,omitempty"`
	CompanyFound bool            `json:"companyFound"`
	Data         *CompanyDetails `json:"data,omitempty"`
}
