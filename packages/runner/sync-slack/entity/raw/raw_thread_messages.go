package rawentity

import (
	"time"
)

type RawThreadMessage struct {
	RawId     string    `gorm:"column:raw_id;default:gen_random_uuid();primaryKey"`
	Data      string    `gorm:"type:jsonb;not null"`
	EmittedAt time.Time `gorm:"type:timestamp with time zone;default:now()"`
}

func (RawThreadMessage) TableName() string {
	return "_openline_raw_thread_messages"
}
