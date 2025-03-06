package postgres_entity

import (
	"github.com/google/uuid"
	"time"
)

type IngestEmailImportStatePeriod string

const (
	REAL_TIME           IngestEmailImportStatePeriod = "REAL_TIME"
	HISTORY             IngestEmailImportStatePeriod = "HISTORY" // this is not used in DB. this is used just in code to trigger the import for the other states
	LAST_WEEK           IngestEmailImportStatePeriod = "LAST_WEEK"
	LAST_3_MONTHS       IngestEmailImportStatePeriod = "LAST_3_MONTHS"
	LAST_YEAR           IngestEmailImportStatePeriod = "LAST_YEAR"
	OLDER_THAN_ONE_YEAR IngestEmailImportStatePeriod = "OLDER_THAN_ONE_YEAR"
)

type IngestEmailImportState struct {
	ID        uuid.UUID                    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant    string                       `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Username  string                       `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Provider  string                       `gorm:"size:255;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	Period    IngestEmailImportStatePeriod `gorm:"size:50;not null;uniqueIndex:uq_one_state_per_tenant_and_user"`
	StartDate *time.Time                   `gorm:""`
	StopDate  *time.Time                   `gorm:""`
	Active    bool                         `gorm:"not null"`
	Cursor    string                       `gorm:"size:255;not null"`
}

func (IngestEmailImportState) TableName() string {
	return "ingest_email_import_state"
}
