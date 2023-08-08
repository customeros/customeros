package rawentity

import (
	"github.com/google/uuid"
	"time"
)

type RawChannelMessage struct {
	ID        uuid.UUID `gorm:"column:openline_raw_id;type:uuid;default:gen_random_uuid();primaryKey"`
	Data      string    `gorm:"type:jsonb;not null"`
	EmittedAt time.Time `gorm:"type:timestamp with time zone;default:now()"`
}

func (RawChannelMessage) TableName() string {
	return "_openline_raw_cahnnel_messages"
}
