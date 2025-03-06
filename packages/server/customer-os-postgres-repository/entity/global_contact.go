package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type GlobalContact struct {
	ID                      uint64              `gorm:"primary_key;autoIncrement" json:"id"`
	CreatedAt               time.Time           `gorm:"type:timestamp;default:current_timestamp" json:"created_at"`
	UpdatedAt               time.Time           `gorm:"type:timestamp;default:current_timestamp" json:"updated_at"`
	FirstName               string              `gorm:"type:varchar(255)" json:"first_name"`
	LastName                string              `gorm:"type:varchar(255)" json:"last_name"`
	LinkedInIdentifier      string              `gorm:"type:varchar(255);index:idx_linkedin_identifier" json:"linkedin_identifier"`
	WorkEmail               string              `gorm:"type:varchar(255);index:idx_work_email" json:"work_email,omitempty"`
	PersonalEmail           string              `gorm:"type:varchar(255);index:idx_personal_email" json:"personal_email,omitempty"`
	JobTitle                string              `gorm:"type:varchar(255)" json:"job_title"`
	JobStartedAt            *time.Time          `gorm:"type:timestamp;null" json:"job_started_at,omitempty"`
	JobEndedAt              *time.Time          `gorm:"type:timestamp;null" json:"job_ended_at,omitempty"`
	PrimaryDomain           string              `gorm:"type:varchar(255);index:idx_primary_domain" json:"primary_domain"`
	ProfilePhotoExternalUrl string              `gorm:"type:varchar(1000)" json:"profile_photo_external_url,omitempty"`
	ProfilePhotoPath        string              `gorm:"column:profile_photo_path;type:varchar(2000)" json:"profilePhotoPath"`
	DownloadStatus          enum.DownloadStatus `gorm:"column:download_status;type:varchar(55);default:'NOT_STARTED'" json:"downloadStatus"`
}

// TableName specifies the table name for the GlobalContact entity
func (GlobalContact) TableName() string {
	return "global_contacts"
}

// Composite indexes for unique combinations ???
// CREATE INDEX idx_linkedin_identifier_primary_domain ON global_contacts (linkedin_identifier, primary_domain);
// CREATE INDEX idx_work_email_primary_domain ON global_contacts (work_email, primary_domain);
