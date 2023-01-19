package entity

import "time"

type Company struct {
	Id                     string    `gorm:"column:id"`
	AirbyteAbId            string    `gorm:"column:_airbyte_ab_id"`
	AirbyteCompaniesHashid string    `gorm:"column:_airbyte_companies_hashid"`
	CreateDate             time.Time `gorm:"column:createdat"`
	UpdatedDate            time.Time `gorm:"column:updatedat"`
}

type Companies []Company

func (Company) TableName() string {
	return "companies"
}
