package postgres_entity

import (
	"time"
)

type GlobalContact struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt          time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	FirstName          string     `gorm:"type:varchar(255)" json:"first_name"`
	LastName           string     `gorm:"type:varchar(255)" json:"last_name"`
	LinkedInIdentifier string     `gorm:"type:varchar(255);index:idx_linkedin_identifier" json:"linkedin_identifier"`
	WorkEmail          string     `gorm:"type:varchar(255);index:idx_work_email" json:"work_email"`
	PersonalEmail      string     `gorm:"type:varchar(255);index:idx_personal_email" json:"personal_email"`
	JobTitle           string     `gorm:"type:varchar(255)" json:"job_title"`
	JobStartedAt       *time.Time `gorm:"type:timestamp;null" json:"job_started_at,omitempty"`
	JobEndedAt         *time.Time `gorm:"type:timestamp;null" json:"job_ended_at,omitempty"`
	PrimaryDomain      string     `gorm:"type:varchar(255);index:idx_primary_domain" json:"primary_domain"`
}

// TableName specifies the table name for the GlobalContact entity
func (GlobalContact) TableName() string {
	return "global_contacts"
}

// Composite indexes for unique combinations ???
// CREATE INDEX idx_linkedin_identifier_primary_domain ON global_contacts (linkedin_identifier, primary_domain);
// CREATE INDEX idx_work_email_primary_domain ON global_contacts (work_email, primary_domain);
