package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
)

type AgentCapabilityRegistry struct {
	ID          string `gorm:"primaryKey;type:varchar(20)" json:"id"`
	Type        string `gorm:"column:type;type:varchar(255)" json:"type" binding:"required"`
	Description string `gorm:"column:description;type:varchar(255);index" json:"description" binding:"required"`
	Action      string `gorm:"column:action;type:varchar(255)" json:"action"`
	IsActive    bool   `gorm:"column:is_active;type:boolean;default:false" json:"isActive"`
}

func (AgentCapabilityRegistry) TableName() string {
	return "agent_capability_registry"
}

func (r *AgentCapabilityRegistry) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("acr", 16)
	return nil
}
