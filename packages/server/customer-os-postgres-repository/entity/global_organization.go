package entity

import "time"

type GlobalOrganization struct {
	ID                          uint64    `gorm:"primary_key;autoIncrement" json:"id"`
	Name                        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	PrimaryDomain               string    `gorm:"column:primary_domain;type:varchar(255);NOT NULL;index:idx_global_organization_primary_domain,unique" json:"primaryDomain"`
	Domains                     []string  `gorm:"column:domains;type:text[]" json:"domains"`
	CreatedAt                   time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt                   time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Description                 string    `gorm:"column:description;type:text" json:"description"`
	IconUrl                     string    `gorm:"column:icon_url;type:varchar(2000)" json:"iconUrl"`
	LogoUrl                     string    `gorm:"column:logo_url;type:varchar(2000)" json:"logoUrl"`
	Website                     string    `gorm:"column:website;type:varchar(255)" json:"website"`
	LinkedInUrl                 string    `gorm:"column:linkedin;type:varchar(255)" json:"linkedin"`
	LinkedInAlias               string    `gorm:"column:linkedin_alias;type:varchar(255)" json:"linkedinAlias"`
	OtherSocials                []string  `gorm:"column:other_socials;type:text[]" json:"otherSocials"`
	ValueProposition            string    `gorm:"column:value_proposition;type:text" json:"valueProposition"`
	TargetAudience              string    `gorm:"column:target_audience;type:text" json:"targetAudience"`
	IndustryGicsSectorId        string    `gorm:"column:industry_gics_sector_id;type:varchar(255)" json:"industryGicsSectorId"`
	IndustryGicsIndustryGroupId string    `gorm:"column:industry_gics_industry_group_id;type:varchar(255)" json:"industryGicsIndustryGroupId"`
	IndustryGicsIndustryId      string    `gorm:"column:industry_gics_industry_id;type:varchar(255)" json:"industryGicsIndustryId"`
	IndustryGicsSubIndustryId   string    `gorm:"column:industry_gics_sub_industry_id;type:varchar(255)" json:"industryGicsSubIndustryId"`
	YearFounded                 int       `gorm:"column:year_founded" json:"yearFounded"`
	EmployeeCount               int64     `gorm:"column:employee_count" json:"employeeCount"`
}

// TableName sets the name of the table for GORM
func (GlobalOrganization) TableName() string {
	return "global_organization"
}

//CREATE EXTENSION IF NOT EXISTS pg_trgm;
//CREATE INDEX idx_global_organization_name_trgm ON global_organization USING gin (name gin_trgm_ops);
//CREATE INDEX idx_global_organization_primary_domain_trgm ON global_organizationgo USING gin (primary_domain gin_trgm_ops);
