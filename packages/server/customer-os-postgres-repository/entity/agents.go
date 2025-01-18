package entity

import (
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type Agents struct {
	ID           string     `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Type         string     `gorm:"column:type;type:varchar(50);not null;" json:"type" binding:"required"`
	Tenant       string     `gorm:"column:tenant;type:varchar(255);not null;uniqueIndex:idx_tenant_name" json:"tenant" binding:"required"`
	Name         string     `gorm:"column:name;type:varchar(255);not null" json:"name" binding:"required"`
	Capabilities string     `gorm:"column:capabilities;type:text" json:"capabilities"`
	Goal         string     `gorm:"column:goal;type:text" json:"goal"`
	Status       string     `gorm:"column:status;type:varchar(32)" json:"status"`
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

func (r *Agents) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("agent", 16)
	return nil
}
