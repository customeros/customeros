package entity

import (
	"github.com/google/uuid"
	"time"
)

type CosApiEnrichPersonTempResult struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant                string    `gorm:"column:tenant;type:varchar(255);" json:"tenant"`
	CreatedAt             time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	ScrapinRecordId       uint64    `gorm:"column:scrapin_record_id;type:bigint;" json:"scrapinRecordId"`
	BettercontactRecordId string    `gorm:"column:bettercontact_record_id;type:varchar(255);" json:"bettercontactRecordId"`
}

func (CosApiEnrichPersonTempResult) TableName() string {
	return "cos_api_enrich_person_temp_result"
}
