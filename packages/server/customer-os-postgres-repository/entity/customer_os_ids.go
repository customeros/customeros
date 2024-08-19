package entity

import "time"

const Organization = "organization"

type CustomerOsIds struct {
	Tenant       string    `gorm:"column:tenant;size:50;primaryKey"`
	CustomerOSID string    `gorm:"column:customer_os_id;size:30;primaryKey"`
	Entity       string    `gorm:"column:entity;size:30"`
	EntityId     string    `gorm:"column:entity_id;size:50"`
	CreatedDate  time.Time `gorm:"default:current_timestamp"`
	Attempts     int       `gorm:"column:attempts"`
}

func (CustomerOsIds) TableName() string {
	return "customer_os_ids"
}
