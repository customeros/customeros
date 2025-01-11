package entity

import (
	"time"
)

type Agents struct {
	ID           string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	RegistryID   string     `gorm:"column:registry_id;type:varchar(50);not null;" json:"registryId" binding:"required"`
	Tenant       string     `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name         string     `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"name" binding:"required"`
	Capabilities string     `gorm:"column:capabilities;type:text" json:"capabilities"`
	Goal         string     `gorm:"column:goal;type:text" json:"goal"`
	Config       *string    `gorm:"column:config;type:text" json:"config"`
	IsActive     bool       `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	FlowID       string     `gorm:"column:flow_id;type:varchar(255)" json:"flowId" binding:"required"`
	VisibleInUI  bool       `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ErrorMessage *string    `gorm:"column:error_message;type:varchar(255)" json:"errorMessage"`
	Color        string     `gorm:"column:color;type:varchar(255)" json:"color"`
	Icon         string     `gorm:"column:icon;type:varchar(255)" json:"icon"`
}

func (Agents) TableName() string {
	return "agents"
}
