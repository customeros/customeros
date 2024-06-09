package entity

import (
	"time"
)

type AirbyteRaw struct {
	AirbyteRawId       string    `gorm:"column:_airbyte_raw_id"`
	AirbyteData        string    `gorm:"column:_airbyte_data"`
	AirbyteExtractedAt time.Time `gorm:"column:_airbyte_extracted_at"`
}

type AirbyteRaws []AirbyteRaw
