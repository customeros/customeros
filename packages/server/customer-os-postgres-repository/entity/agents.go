package entity

import (
	"time"
)

type Agents struct {
	ID             string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	RegistryID     string     `gorm:"column:registry_id;type:varchar(50);not null;" json:"registryId" binding:"required"`
	Tenant         string     `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name           string     `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"name" binding:"required"`
	Description    *string    `gorm:"column:description;type:text" json:"description"`
	IsActive       bool       `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	FlowID         string     `gorm:"column:flow_id;type:varchar(255)" json:"flowId" binding:"required"`
	VisibleInUI    *bool      `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	UserContext    *string    `gorm:"column:user_context;type:text" json:"input"`
	Config         *string    `gorm:"column:config;type:text" json:"config"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy      *string    `gorm:"column:created_by;type:varchar(255)" json:"createdBy"`
	LastModifiedBy *string    `gorm:"column:last_modified_by;type:varchar(255)" json:"lastModifiedBy"`
}

func (Agents) TableName() string {
	return "agents"
}
