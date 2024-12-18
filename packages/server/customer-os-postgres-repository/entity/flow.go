package entity

import (
	"time"
)

type Flow struct {
	ID             string     `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Tenant         string     `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name           string     `gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"name" binding:"required"`
	Description    *string    `gorm:"column:description;type:text" json:"description"`
	TriggerOn      string     `gorm:"column:trigger_on;type:varchar(255);not null" json:"triggerOn" binding:"required"`
	TriggerNodeID  string     `gorm:"column:trigger_node_id;type:varchar(255);not null" json:"triggerNodeId" binding:"required"`
	VisibleInUI    *bool      `gorm:"column:visible_in_ui;type:boolean;default:true" json:"visibleInUI"`
	Status         string     `gorm:"column:status;type:varchar(50);not null;default:'inactive'" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CreatedBy      *string    `gorm:"column:created_by;type:varchar(255)" json:"createdBy"`
	LastModifiedBy *string    `gorm:"column:last_modified_by;type:varchar(255)" json:"lastModifiedBy"`
}

func (Flow) TableName() string {
	return "flow"
}
