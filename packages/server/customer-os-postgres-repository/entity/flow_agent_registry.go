package entity

import "gorm.io/datatypes"

type FlowAgentRegistry struct {
	ID           uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Agent        string          `gorm:"column:agent;type:varchar(255);not null;index" json:"agent" binding:"required"`
	FriendlyName string          `gorm:"column:friendly_name;type:varchar(255);not null;index" json:"FriendlyName" binding:"required"`
	Description  string          `gorm:"column:description;type:varchar(255);not null" json:"description" binding:"required"`
	Instructions *datatypes.JSON `gorm:"column:instructions;type:jsonb" json:"instructions"`
	Status       string          `gorm:"column:status;type:varchar(50)" json:"status"`
}

func (FlowAgentRegistry) TableName() string {
	return "flow_agent_registry"
}
