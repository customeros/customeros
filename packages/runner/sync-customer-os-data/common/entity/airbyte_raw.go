package entity

import (
	"time"
)

type AirbyteRaw struct {
	AirbyteAbId      string    `gorm:"column:_airbyte_ab_id"`
	AirbyteData      string    `gorm:"column:_airbyte_data"`
	AirbyteEmittedAt time.Time `gorm:"column:_airbyte_emitted_at"`
}

type AirbyteRaws []AirbyteRaw
