package entity

import "time"

type CompanyRaw struct {
	AirbyteAbId      string    `gorm:"column:_airbyte_ab_id"`
	AirbyteData      string    `gorm:"column:_airbyte_data"`
	AirbyteEmittedAt time.Time `gorm:"column:_airbyte_emitted_at"`
}

type CompaniesRaw []CompanyRaw

func (CompanyRaw) TableName() string {
	return "_airbyte_raw_companies"
}
