package entity

import (
	"time"
)

type EnrichDetailsBetterContact struct {
	ID                 string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RequestID          string    `gorm:"column:request_id;type:varchar(255);NOT NULL" json:"requestId"`
	CreatedAt          time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	ContactFirstName   string    `gorm:"column:contact_first_name;type:varchar(255);" json:"contactFirstName"`
	ContactLastName    string    `gorm:"column:contact_last_name;type:varchar(255);" json:"contactLastName"`
	ContactLinkedInUrl string    `gorm:"column:contact_linkedin_url;type:varchar(255);" json:"contactLinkedInUrl"`
	CompanyName        string    `gorm:"column:company_name;type:varchar(255);" json:"companyName"`
	CompanyDomain      string    `gorm:"column:company_domain;type:varchar(255);" json:"companyDomain"`
	EnrichPhoneNumber  bool      `gorm:"column:enrich_phone_number;type:boolean;DEFAULT:false" json:"enrichPhoneNumber"`
	Request            string    `gorm:"column:request;type:text;" json:"request"`
	Response           string    `gorm:"column:response;type:text;" json:"response"`
}

func (EnrichDetailsBetterContact) TableName() string {
	return "enrich_details_better_contact"
}

type BetterContactResponseBody struct {
	Id      string `json:"id"`
	Status  string `json:"status"`
	Summary struct {
		Total         int `json:"total"`
		Valid         int `json:"valid"`
		Undeliverable int `json:"undeliverable"`
		NotFound      int `json:"not_found"`
	} `json:"summary"`
	Data []struct {
		Enriched                     bool          `json:"enriched"`
		EmailProvider                string        `json:"email_provider"`
		CompanyName                  string        `json:"company_name"`
		CompanyDomain                string        `json:"company_domain"`
		ContactGender                interface{}   `json:"contact_gender"`
		CompanyIndustry              interface{}   `json:"company_industry"`
		CompanyLegalId               interface{}   `json:"company_legal_id"`
		ContactJobTitle              interface{}   `json:"contact_job_title"`
		ContactLastName              string        `json:"contact_last_name"`
		CompanyLegalName             interface{}   `json:"company_legal_name"`
		ContactFirstName             string        `json:"contact_first_name"`
		CompanyAddressCity           interface{}   `json:"company_address_city"`
		CompanyLinkedinUrl           interface{}   `json:"company_linkedin_url"`
		CompanyPhoneNumber           interface{}   `json:"company_phone_number"`
		ContactPhoneNumber           interface{}   `json:"contact_phone_number"`
		CompanyAddressState          interface{}   `json:"company_address_state"`
		CompanyIndustryCode          interface{}   `json:"company_industry_code"`
		ContactEmailAddress          string        `json:"contact_email_address"`
		CompanyAddressStreet         interface{}   `json:"company_address_street"`
		CompanyAddressCountry        interface{}   `json:"company_address_country"`
		CompanyAddressZipcode        interface{}   `json:"company_address_zipcode"`
		CompanyEmployeesNumber       interface{}   `json:"company_employees_number"`
		ContactEmailAddressStatus    string        `json:"contact_email_address_status"`
		ContactLinkedinProfileUrl    string        `json:"contact_linkedin_profile_url"`
		ContactEmailAddressProvider  string        `json:"contact_email_address_provider"`
		ContactAdditionalPhoneNumber interface{}   `json:"contact_additional_phone_number"`
		CustomFields                 []interface{} `json:"custom_fields"`
	} `json:"data"`
}
