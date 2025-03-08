package postgres_entity

import (
	"time"
)

// CacheCrustData represents the cache table for Crust Data API responses
type CacheCrustData struct {
	ID                   uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RequestCompanyDomain string    `gorm:"column:request_company_domain;type:varchar(255)" json:"requestCompany"` // Organization primary domain
	RequestJobTitle      string    `gorm:"column:request_job_title;varchar(255)" json:"requestJobTitle"`          // List of job titles
	Response             string    `gorm:"column:response;type:text" json:"response"`
	CreatedAt            time.Time `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt            time.Time `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
}

func (CacheCrustData) TableName() string {
	return "cache_crust_data"
}
