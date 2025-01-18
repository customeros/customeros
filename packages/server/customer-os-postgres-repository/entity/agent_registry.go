package entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type AgentRegistry struct {
	ID           string `gorm:"primaryKey;type:varchar(20)" json:"id"`
	Type         string `gorm:"column:type;type:varchar(255)" json:"type" binding:"required"`
	Name         string `gorm:"column:name;type:varchar(255);not null;index" json:"name" binding:"required"`
	Capabilities string `gorm:"column:capabilities;type:text" json:"capabilities"`
	Goal         string `gorm:"column:goal;type:text" json:"goal"`
	IsActive     bool   `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
	Icon         string `gorm:"column:icon;type:text" json:"icon"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}

func (r *AgentRegistry) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("ar", 16)
	return nil
}

