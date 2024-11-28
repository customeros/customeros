package entity

import "time"

type Tenant struct {
	Name      string    `gorm:"primary_key;type:varchar(255);NOT NULL" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	HashID    string    `gorm:"column:tenant_hash;type:varchar(25)" json:"hashID"`
}

func (Tenant) TableName() string {
	return "tenant"
}
