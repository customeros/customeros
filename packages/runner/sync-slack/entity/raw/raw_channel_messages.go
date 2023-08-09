package rawentity

import (
	"time"
)

type RawChannelMessage struct {
	RawId     string    `gorm:"column:raw_id;default:gen_random_uuid();primaryKey"`
	Data      string    `gorm:"type:jsonb;not null"`
	EmittedAt time.Time `gorm:"type:timestamp with time zone;default:now()"`
}

func (RawChannelMessage) TableName() string {
	return "_openline_raw_channel_messages"
}
