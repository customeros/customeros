package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type EnrichDetailsTracking struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	IP             string    `gorm:"column:ip;uniqueIndex:ip_unique;type:varchar(255);" json:"ip"`
	CompanyName    *string   `gorm:"column:company_name;type:varchar(255);" json:"companyName"`
	CompanyDomain  *string   `gorm:"column:company_domain;type:varchar(255);" json:"companyDomain"`
	CompanyWebsite *string   `gorm:"column:company_website;type:varchar(255);" json:"companyWebsite"`
	Response       string    `gorm:"column:response;type:text;" json:"response"`
}

func (EnrichDetailsTracking) TableName() string {
	return "enrich_details_tracking"
}

type SnitcherResponseBody struct {
	Fuzzy   bool   `json:"fuzzy"`
	Domain  string `json:"domain"`
	Type    string `json:"type"`
	Company *struct {
		Name          string      `json:"name"`
		Domain        string      `json:"domain"`
		Website       string      `json:"website"`
		Industry      string      `json:"industry"`
		FoundedYear   interface{} `json:"founded_year"`
		EmployeeRange string      `json:"employee_range"`
		AnnualRevenue interface{} `json:"annual_revenue"`
		TotalFunding  interface{} `json:"total_funding"`
		Location      Location    `json:"location"`
		Description   string      `json:"description"`
		Phone         string      `json:"phone"`
		Geo           struct {
			Country      string  `json:"country"`
			CountryCode  string  `json:"country_code"`
			State        string  `json:"state"`
			StateCode    *string `json:"state_code"`
			PostalCode   *string `json:"postal_code"`
			City         string  `json:"city"`
			Street       *string `json:"street"`
			StreetNumber *string `json:"street_number"`
		} `json:"geo"`
		Profiles *struct {
			Crunchbase *struct {
				Handle string      `json:"handle"`
				Url    interface{} `json:"url"`
			} `json:"crunchbase"`
			Linkedin *struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"linkedin"`
			Facebook *struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"facebook"`
			Twitter *struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"twitter"`
			Instagram *struct {
				Handle string      `json:"handle"`
				Url    interface{} `json:"url"`
			} `json:"instagram"`
			Youtube *struct {
				Handle string      `json:"handle"`
				Url    interface{} `json:"url"`
			} `json:"youtube"`
		} `json:"profiles"`
	} `json:"company"`
	GeoIP struct {
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
		City        string `json:"city"`
		State       string `json:"state"`
	} `json:"geoIP"`
}

// Define a custom type for Location
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

func (location SnitcherResponseBody) HasLocation() bool {
	return location.Company.Location.CityName != "" || location.Company.Location.Country.Name != "" || location.Company.Location.RegionName != "" || location.Company.Location.Country.Iso2 != ""
}

func (location SnitcherResponseBody) HasRawLocation() bool {
	return location.Company.Location.RawLocation != ""
}

func (location SnitcherResponseBody) LocationToString() string {
	result := ""
	if location.Company.Location.StreetName != "" {
		result += location.Company.Location.StreetName
	}
	if location.Company.Location.StreetNumber != "" {
		result += " " + location.Company.Location.StreetNumber
	}
	if location.Company.Location.PostalCode != "" {
		if result != "" {
			result += ", "
		}
		result += location.Company.Location.PostalCode
	}
	if location.Company.Location.CityName != "" {
		if result != "" {
			result += ", "
		}
		result += location.Company.Location.CityName
	}
	if location.Company.Location.RegionName != "" {
		if result != "" {
			result += ", "
		}
		result += location.Company.Location.RegionName
	}
	if location.Company.Location.Country.Name != "" {
		if result != "" {
			result += ", "
		}
		result += location.Company.Location.Country.Name
	}

	if result == "" && location.HasRawLocation() {
		result = location.Company.Location.RawLocation
	}

	if result == "" {
		return "Unknown"
	}

	return result
}
