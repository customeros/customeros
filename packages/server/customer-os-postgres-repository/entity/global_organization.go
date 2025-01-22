package postgres_entity

import "time"

type GlobalOrganization struct {
	ID                      uint64     `gorm:"primary_key;autoIncrement" json:"id"`
	Name                    string     `gorm:"column:name;type:varchar(255)" json:"name"`
	PrimaryDomain           string     `gorm:"column:primary_domain;type:varchar(255);NOT NULL;index:idx_global_organization_primary_domain,unique" json:"primaryDomain"`
	OtherDomains            string     `gorm:"column:other_domains;type:text" json:"otherDomains"`
	CreatedAt               time.Time  `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt               time.Time  `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Description             string     `gorm:"column:description;type:text" json:"description"`
	IconUrl                 string     `gorm:"column:icon_url;type:varchar(2000)" json:"iconUrl"`
	LogoUrl                 string     `gorm:"column:logo_url;type:varchar(2000)" json:"logoUrl"`
	Website                 string     `gorm:"column:website;type:varchar(255)" json:"website"`
	LinkedInUrl             string     `gorm:"column:linkedin;type:varchar(255)" json:"linkedin"`
	LinkedInAlias           string     `gorm:"column:linkedin_alias;type:varchar(255)" json:"linkedinAlias"`
	OtherSocials            []string   `gorm:"column:other_socials;type:text[]" json:"otherSocials"`
	IndustryNaicsCode       string     `gorm:"column:industry_naics_code;type:varchar(255)" json:"industryNaicsCode"`
	IndustryNaicsName       string     `gorm:"column:industry_naics_name;type:varchar(255)" json:"industryNaicsName"`
	IndustrySetAt           *time.Time `gorm:"column:industry_set_at;type:timestamp" json:"industrySetAt"`
	IndustryRequestedAt     *time.Time `gorm:"column:industry_requested_at;type:timestamp" json:"industryRequestedAt"`
	IndustryRequestCount    int        `gorm:"column:industry_request_count" json:"industryRequestCount"`
	Market                  string     `gorm:"column:market;type:text" json:"market"`
	YearFounded             int        `gorm:"column:year_founded" json:"yearFounded"`
	EmployeeCount           int64      `gorm:"column:employee_count" json:"employeeCount"`
	City                    string     `gorm:"column:city;type:varchar(255)" json:"city"`
	Region                  string     `gorm:"column:region;type:varchar(255)" json:"region"`
	CountryA2               string     `gorm:"column:country_a2;type:varchar(255)" json:"countryA2"`
	DescriptionSetAt        *time.Time `gorm:"column:description_set_at;type:timestamp" json:"descriptionSetAt"`
	DescriptionRequestedAt  *time.Time `gorm:"column:description_requested_at;type:timestamp" json:"descriptionRequestedAt"`
	DescriptionRequestCount int        `gorm:"column:description_request_count" json:"descriptionRequestCount"`
	NameSetAt               *time.Time `gorm:"column:name_set_at;type:timestamp" json:"nameSetAt"`
	NameRequestedAt         *time.Time `gorm:"column:name_requested_at;type:timestamp" json:"nameRequestedAt"`
	NameRequestCount        int        `gorm:"column:name_request_count" json:"nameRequestCount"`
	SourceDescription1      string     `gorm:"column:source_description_1;type:text" json:"sourceDescription1"`
	SourceDescription2      string     `gorm:"column:source_description_2;type:text" json:"sourceDescription2"`
	SourceDescription3      string     `gorm:"column:source_description_3;type:text" json:"sourceDescription3"`
	SourceDescription4      string     `gorm:"column:source_description_4;type:text" json:"sourceDescription4"`
	SourceDescription5      string     `gorm:"column:source_description_5;type:text" json:"sourceDescription5"`
	SyncedToNeoAt           *time.Time `gorm:"column:synced_to_neo_at;type:timestamp" json:"syncedToNeoAt"`
}

// TableName sets the name of the table for GORM
func (GlobalOrganization) TableName() string {
	return "global_organization"
}

//CREATE EXTENSION IF NOT EXISTS pg_trgm;
//CREATE INDEX idx_global_organization_name_trgm ON global_organization USING gin (name gin_trgm_ops);
//CREATE INDEX idx_global_organization_primary_domain_trgm ON global_organization USING gin (primary_domain gin_trgm_ops);
//CREATE INDEX idx_global_organization_other_domains_trgm ON global_organization USING gin (other_domains gin_trgm_ops);
