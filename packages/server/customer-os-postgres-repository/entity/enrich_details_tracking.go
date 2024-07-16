package entity

import (
	"time"
)

type EnrichDetailsTracking struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	IP            string    `gorm:"column:ip;uniqueIndex:ip_unique;type:varchar(255);" json:"ip"`
	CompanyName   *string   `gorm:"column:company_name;type:varchar(255);" json:"companyName"`
	CompanyDomain *string   `gorm:"column:company_domain;type:varchar(255);" json:"companyDomain"`
	Response      string    `gorm:"column:response;type:text;" json:"response"`
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
		FoundedYear   string      `json:"founded_year"`
		EmployeeRange string      `json:"employee_range"`
		AnnualRevenue interface{} `json:"annual_revenue"`
		TotalFunding  interface{} `json:"total_funding"`
		Location      string      `json:"location"`
		Description   string      `json:"description"`
		Phone         string      `json:"phone"`
		Geo           struct {
			Country      string `json:"country"`
			CountryCode  string `json:"country_code"`
			State        string `json:"state"`
			StateCode    string `json:"state_code"`
			PostalCode   string `json:"postal_code"`
			City         string `json:"city"`
			Street       string `json:"street"`
			StreetNumber string `json:"street_number"`
		} `json:"geo"`
		Profiles struct {
			Crunchbase struct {
				Handle string      `json:"handle"`
				Url    interface{} `json:"url"`
			} `json:"crunchbase"`
			Linkedin struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"linkedin"`
			Facebook struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"facebook"`
			Twitter struct {
				Handle string `json:"handle"`
				Url    string `json:"url"`
			} `json:"twitter"`
			Instagram struct {
				Handle string      `json:"handle"`
				Url    interface{} `json:"url"`
			} `json:"instagram"`
			Youtube struct {
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
