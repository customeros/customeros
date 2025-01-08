package entity

import (
	"time"
)

type Automations struct {
	ID             string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	RegistryID     string     `gorm:"column:registry_id;type:varchar(50);not null;" json:"registryId" binding:"required"`
	Tenant         string     `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name           string     `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"name" binding:"required"`
	Description    *string    `gorm:"column:description;type:text" json:"description"`
	IsActive       string     `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	ListensFor     string     `gorm:"column:listens_for;type:varchar(255);not null" json:"listensFor" binding:"required"`
	AgentID        string     `gorm:"column:agent_id;type:varchar(255)" json:"agentId" binding:"required"`
	FlowID         string     `gorm:"column:flow_id;type:varchar(255)" json:"flowId" binding:"required"`
	VisibleInUI    *bool      `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy      *string    `gorm:"column:created_by;type:varchar(255)" json:"createdBy"`
	LastModifiedBy *string    `gorm:"column:last_modified_by;type:varchar(255)" json:"lastModifiedBy"`
}

func (Automations) TableName() string {
	return "automations"
}
