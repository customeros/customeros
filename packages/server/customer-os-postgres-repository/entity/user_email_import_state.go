package entity

import (
	"github.com/google/uuid"
	"time"
)

type EmailImportState string

const (
	REAL_TIME           EmailImportState = "REAL_TIME"
	HISTORY             EmailImportState = "HISTORY" // this is not used in DB. this is used just in code to trigger the import for the other states
	LAST_WEEK           EmailImportState = "LAST_WEEK"
	LAST_3_MONTHS       EmailImportState = "LAST_3_MONTHS"
	LAST_YEAR           EmailImportState = "LAST_YEAR"
	OLDER_THAN_ONE_YEAR EmailImportState = "OLDER_THAN_ONE_YEAR"
)

type UserEmailImportState struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant    string           `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Username  string           `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Provider  string           `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	State     EmailImportState `gorm:"size:50;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	StartDate *time.Time       `gorm:""`
	StopDate  *time.Time       `gorm:""`
	Active    bool             `gorm:"not null"`
	Cursor    string           `gorm:"size:255;not null"`
}

func (UserEmailImportState) TableName() string {
	return "user_email_import_state"
}
