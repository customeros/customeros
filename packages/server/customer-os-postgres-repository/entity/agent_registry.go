package postgres_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AgentRegistry struct {
	ID           string         `gorm:"primaryKey;type:varchar(21)" json:"id"`
	Type         enum.AgentType `gorm:"column:type;type:varchar(255);not null;index" json:"type" binding:"required"`
	AgentName    string         `gorm:"column:agent_name;type:varchar(255);not null" json:"agentName"`
	Filename     string         `gorm:"column:filename;type:varchar(255);not null" json:"filename"`
	Triggers     pq.StringArray `gorm:"column:triggers;type:varchar[]" json:"triggers"`
	Capabilities pq.StringArray `gorm:"column:capabilities;type:varchar[]" json:"capabilities"`
	Version      string         `gorm:"column:version;type:varchar(21);not null" json:"version"`
	Goal         string         `gorm:"column:goal;type:varchar(255)" json:"goal"`
	IsActive     bool           `gorm:"column:is_active;type:boolean;default:true" json:"isActive"`
	Icon         string         `gorm:"column:icon;type:text" json:"icon"`
}

func (AgentRegistry) TableName() string {
	return "agent_registry"
}

func (r *AgentRegistry) BeforeCreate(tx *gorm.DB) error {
	r.ID = utils.GenerateNanoIdWithPrefix("ar", 16)
	return nil
}
