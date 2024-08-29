package entity

import (
	"time"
)

type EnrichDetailsPreFilterTracking struct {
	ID                 string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt          time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:timestamp;" json:"updatedAt"`
	IP                 string    `gorm:"column:ip;uniqueIndex:ip_unique;type:varchar(255);" json:"ip"`
	ShouldIdentify     *bool     `gorm:"column:should_identify;type:boolean;" json:"shouldIdentify"`
	SkipIdentifyReason string    `gorm:"column:skip_identify_reason;type:varchar(255);" json:"skipIdentifyReason"`
	Response           *string   `gorm:"column:response;type:text;" json:"response"`
}

func (EnrichDetailsPreFilterTracking) TableName() string {
	return "enrich_details_prefilter_tracking"
}
