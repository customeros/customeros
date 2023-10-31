package entity

import (
	"github.com/google/uuid"
	"time"
)

type GmailImportState string

const (
	REAL_TIME           GmailImportState = "REAL_TIME"
	HISTORY             GmailImportState = "HISTORY" // this is not used in DB. this is used just in code to trigger the import for the other states
	LAST_WEEK           GmailImportState = "LAST_WEEK"
	LAST_3_MONTHS       GmailImportState = "LAST_3_MONTHS"
	LAST_YEAR           GmailImportState = "LAST_YEAR"
	OLDER_THAN_ONE_YEAR GmailImportState = "OLDER_THAN_ONE_YEAR"
)

type UserGmailImportState struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant    string           `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Username  string           `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	State     GmailImportState `gorm:"size:50;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	StartDate *time.Time       `gorm:""`
	StopDate  *time.Time       `gorm:""`
	Active    bool             `gorm:"not null"`
	Cursor    string           `gorm:"size:255;not null"`
}

func (UserGmailImportState) TableName() string {
	return "user_gmail_import_state"
}
