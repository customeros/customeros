package models

import "time"

type Contact struct {
	ID              string `gorm:"column:id;primaryKey;type:varchar(30);"`
	Tenant          string `gorm:"column:tenant;type:varchar(255);index;not null"`
	GlobalContactID string `gorm:"column:tenant;type:varchar(255);index;not null"`

	Hide     string    `gorm:"column:tenant;type:varchar(255);index;not null"`
	HiddenAt time.Time `gorm:"column:created_at;autoCreateTime"`

	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
