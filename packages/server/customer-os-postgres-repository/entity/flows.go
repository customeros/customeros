package entity

import (
	"time"
)

type Flows struct {
	ID             string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	AutomationID   string     `gorm:"primaryKey;type:varchar(50)" json:"automationId"`
	IsActive       string     `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	LastModifiedBy *string    `gorm:"column:last_modified_by;type:varchar(255)" json:"lastModifiedBy"`
}

func (Flows) TableName() string {
	return "flows"
}
