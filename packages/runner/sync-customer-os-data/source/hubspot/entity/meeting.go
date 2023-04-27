package entity

import (
	"github.com/jackc/pgtype"
	"time"
)

type Meeting struct {
	Id                    string       `gorm:"column:id"`
	AirbyteAbId           string       `gorm:"column:_airbyte_ab_id"`
	AirbyteMeetingsHashid string       `gorm:"column:_airbyte_engagements_meetings_hashid"`
	CreateDate            time.Time    `gorm:"column:createdat"`
	UpdatedDate           time.Time    `gorm:"column:updatedat"`
	ContactsExternalIds   pgtype.JSONB `gorm:"column:contacts;type:jsonb"`
}

type Meetings []Meeting

func (Meeting) TableName() string {
	return "engagements_meetings"
}
