package entity

import "time"

type OwnerRaw struct {
	AirbyteAbId      string    `gorm:"column:_airbyte_ab_id"`
	AirbyteData      string    `gorm:"column:_airbyte_data"`
	AirbyteEmittedAt time.Time `gorm:"column:_airbyte_emitted_at"`
}

type OwnersRaw []OwnerRaw

func (OwnerRaw) TableName() string {
	return "_airbyte_raw_owners"
}
