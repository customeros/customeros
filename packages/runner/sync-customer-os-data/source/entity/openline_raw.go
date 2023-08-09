package entity

import (
	"time"
)

type OpenlineRaw struct {
	RawId     string    `gorm:"column:raw_id"`
	Data      string    `gorm:"column:data"`
	EmittedAt time.Time `gorm:"column:emitted_at"`
}

type OpenlineRaws []OpenlineRaw
