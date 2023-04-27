package entity

import (
	"database/sql"
	"time"
)

type MeetingProperties struct {
	AirbyteAbId           string          `gorm:"column:_airbyte_ab_id"`
	AirbyteMeetingsHashid string          `gorm:"column:_airbyte_engagements_meetings_hashid"`
	Title                 string          `gorm:"column:hs_meeting_title"`
	CreatedByUserId       sql.NullFloat64 `gorm:"column:hs_created_by_user_id"`
	StartedAt             time.Time       `gorm:"column:hs_meeting_start_time"`
	EndedAt               time.Time       `gorm:"column:hs_meeting_end_time"`
	MeetingExternalUrl    string          `gorm:"column:hs_meeting_external_url"`
	Location              string          `gorm:"column:hs_meeting_location"`
	MeetingHtml           string          `gorm:"column:hs_meeting_body"`
	MeetingText           string          `gorm:"column:hs_body_preview"`
}

type MeetingPropertiesList []EmailProperties

func (MeetingProperties) TableName() string {
	return "engagements_meetings_properties"
}
